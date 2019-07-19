# Copyright (c) 2019 Intel Corporation.
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.
"""EIS Message Envelope utility functions
"""

# Python imports
import json
import warnings
from .exc import MessageBusError

# Cython imports
from .libeismsgbus cimport *
from cpython cimport bool


cdef void put_bytes_helper(msg_envelope_t* env, data):
    """Helper function to serialize a Python bytes object to a blob object.
    """
    cdef msgbus_ret_t ret
    cdef msg_envelope_elem_body_t* body

    body = msgbus_msg_envelope_new_blob(<char*> data, len(data))
    if body == NULL:
        raise MessageBusError('Failed to initialize new blob')

    # Put element into the message envelope
    ret = msgbus_msg_envelope_put(env, NULL, body)
    if ret != msgbus_ret_t.MSG_SUCCESS:
        msgbus_msg_envelope_elem_destroy(body)
        raise MessageBusError('Failed to put blob in message envelope')

    # Setting blob inside of the message envelope to not be owned by the
    # underlying pointer object, because it is owned by the Python interpreter
    # and should not be freed by the message bus
    env.blob.body.blob.shared.owned = <bint> False


cdef msg_envelope_t* python_to_msg_envelope(data):
    """Helper function to create a msg_envelope_t from a Python bytes or
    dictionary object.

    :param data: Data for the message envelope
    :type: bytes or dict
    :return: Message envelope
    :rtype: msg_envelope_t
    """
    cdef msgbus_ret_t ret
    cdef msg_envelope_elem_body_t* body
    cdef msg_envelope_t* env
    cdef content_type_t ct
    cdef char* key = NULL

    binary = None
    kv_data = None

    if isinstance(data, bytes):
        ct = content_type_t.CT_BLOB
        binary = data
    elif isinstance(data, dict):
        ct = content_type_t.CT_JSON
        kv_data = data
    elif isinstance(data, (list, tuple,)):
        ct = content_type_t.CT_JSON
        if len(data) > 2:
            raise MessageBusError('List can only be 2 elements for a msg')

        if isinstance(data[0], bytes):
            if not isinstance(data[1], dict):
                raise MessageBusError('Second element must be dict')

            binary = data[0]
            kv_data = data[1]
        elif isinstance(data[0], dict):
            if not isinstance(data[1], bytes):
                raise MessageBusError('Second element must be bytes')

            binary = data[1]
            kv_data = data[0]
        else:
            raise MessageBusError(
                    f'Unknown data type: {type(data)}, must be bytes or dict')
    else:
        raise MessageBusError(
                'Unable to create msg envelope from type: {}'.format(
                    type(data)))

    # Initialize message envelope object
    env = msgbus_msg_envelope_new(ct)

    if env == NULL:
        raise MessageBusError('Failed to initialize message envelope')

    if ct == content_type_t.CT_BLOB:
        try:
            put_bytes_helper(env, data)
        except MessageBusError:
            msgbus_msg_envelope_destroy(env)
            raise  # Re-raise
    else:
        if binary is not None:
            try:
                put_bytes_helper(env, binary)
            except:
                msgbus_msg_envelope_destroy(env)
                raise  # Re-raise

        for k,v in kv_data.items():
            if isinstance(v, str):
                bv = bytes(v, 'utf-8')
                body = msgbus_msg_envelope_new_string(bv)
            elif isinstance(v, int):
                body = msgbus_msg_envelope_new_integer(<int> v)
            elif isinstance(v, float):
                body = msgbus_msg_envelope_new_floating(<double> v)
            elif isinstance(v, bool):
                body = msgbus_msg_envelope_new_bool(<bint> v)
            else:
                msgbus_msg_envelope_destroy(env)
                raise MessageBusError(f'Unknown data type in dict: {type(v)}')

            k = bytes(k, 'utf-8')
            ret = msgbus_msg_envelope_put(env, <char*> k, body)
            if ret != msgbus_ret_t.MSG_SUCCESS:
                msgbus_msg_envelope_elem_destroy(body)
                msgbus_msg_envelope_destroy(env)
                raise MessageBusError(f'Failed to put element {k}')
            else:
                # The message envelope takes ownership of the memory allocated
                # for these elements. Setting to NULL to keep the state clean.
                body = NULL
                key = NULL

    return env


cdef object msg_envelope_to_python(msg_envelope_t* msg):
    """Convert msg_envelope_t to Python dictionary or bytes object.

    :param msg: Message envelope to convert
    :type: msg_envelope_t*
    """
    cdef msg_envelope_serialized_part_t* parts = NULL

    num_parts = msgbus_msg_envelope_serialize(msg, &parts)
    if num_parts <= 0:
        raise MessageBusError('Error serializing to Python representation')

    if num_parts > 2:
        warnings.warn('The Python library only supports 2 parts!')

    try:
        data = None

        if msg.content_type == content_type_t.CT_JSON:
            data = json.loads(<bytes> parts[0].bytes)
            if num_parts > 1:
                data = (data, <bytes> parts[1].bytes)
        elif msg.content_type == content_type_t.CT_BLOB:
            data = <bytes> parts[0].bytes
        else:
            raise MessageBusError('Unknown content type')

        return data
    finally:
        msgbus_msg_envelope_serialize_destroy(parts, num_parts)
