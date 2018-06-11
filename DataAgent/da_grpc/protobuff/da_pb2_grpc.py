# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
import grpc

import da_pb2 as da__pb2


class daStub(object):
  """The Data Agent service definition.
  """

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.GetConfigInt = channel.unary_unary(
        '/DataAgent.da/GetConfigInt',
        request_serializer=da__pb2.ConfigIntReq.SerializeToString,
        response_deserializer=da__pb2.ConfigIntResp.FromString,
        )
    self.Config = channel.unary_unary(
        '/DataAgent.da/Config',
        request_serializer=da__pb2.ConfigReq.SerializeToString,
        response_deserializer=da__pb2.ConfigResp.FromString,
        )
    self.Query = channel.unary_unary(
        '/DataAgent.da/Query',
        request_serializer=da__pb2.QueryReq.SerializeToString,
        response_deserializer=da__pb2.QueryResp.FromString,
        )


class daServicer(object):
  """The Data Agent service definition.
  """

  def GetConfigInt(self, request, context):
    """**********Internal Interfaces***************
    GetConfigInt internal interface
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def Config(self, request, context):
    """**********External Interfaces***************
    GetConfig external interface
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def Query(self, request, context):
    """GetQuery external interface
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_daServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'GetConfigInt': grpc.unary_unary_rpc_method_handler(
          servicer.GetConfigInt,
          request_deserializer=da__pb2.ConfigIntReq.FromString,
          response_serializer=da__pb2.ConfigIntResp.SerializeToString,
      ),
      'Config': grpc.unary_unary_rpc_method_handler(
          servicer.Config,
          request_deserializer=da__pb2.ConfigReq.FromString,
          response_serializer=da__pb2.ConfigResp.SerializeToString,
      ),
      'Query': grpc.unary_unary_rpc_method_handler(
          servicer.Query,
          request_deserializer=da__pb2.QueryReq.FromString,
          response_serializer=da__pb2.QueryResp.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'DataAgent.da', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))
