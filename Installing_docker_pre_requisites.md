## Manual steps to install docker daemon and docker-compose with proxy settings configuration

1. Install latest docker cli/docker daemon by following    
   https://docs.docker.com/install/linux/docker-ce/ubuntu/#install-docker-ce. Follow `Install using
   the repository` and `Install Docker CE` (follow first 2 steps) sections there. Also, follow the
   manage docker as a non-root user section at   
   https://docs.docker.com/install/linux/linux-postinstall/ to run docker without sudo

2. Please follow the below steps only if the node/system on which the docker setup is tried out is   
   running behind a HTTP proxy server. If that's not the case, this step can be skipped.

   * Configure proxy settings for docker client to connect to internet and for containers to access
     internet by following https://docs.docker.com/network/proxy/. Sample proxy 
     config that could be going into ~/.docker/config.json would look like below
     with appropriate proxy server and ports added.
   
        ```json
        {
            "proxies":
                {
                    "default":
                    {
                        "httpProxy": "http://<proxy_server>:<proxy_port>",
                        "httpsProxy": "http://<proxy_server:<proxy_port>",
                        "noProxy": "127.0.0.1,localhost,<localadmin>"
                    }
                }
        }
        ```

   * Configure proxy settings for docker daemon by following the steps at 
     https://docs.docker.com/config/daemon/systemd/#httphttps-proxy. Use the values for http proxy 
     and https proxy as used in previous step.

     The correct DNS servers need to be updated to the /etc/resolv.conf
     
     ```
     A. Ubuntu 16.04 and earlier

        For Ubuntu 16.04 and earlier, /etc/resolv.conf was dynamically generated by NetworkManager.

        Comment out the line dns=dnsmasq (with a #) in /etc/NetworkManager/NetworkManager.conf

        Restart the NetworkManager to regenerate /etc/resolv.conf :
        sudo systemctl restart network-manager

        Verify on the host: cat /etc/resolv.conf

     B. Ubuntu 18.04 and later

        Ubuntu 18.04 changed to use systemd-resolved to generate /etc/resolv.conf. Now by default it uses a local DNS cache 127.0.0.53. That will not work inside a container, so Docker will default to Google's 8.8.8.8 DNS server, which may break for people behind a firewall.

        /etc/resolv.conf is actually a symlink (ls -l /etc/resolv.conf) which points to /run/systemd/resolve/stub-resolv.conf (127.0.0.53) by default in Ubuntu 18.04.

        Just change the symlink to point to /run/systemd/resolve/resolv.conf, which lists the real DNS servers:
        sudo ln -sf /run/systemd/resolve/resolv.conf /etc/resolv.conf

        Verify on the host: cat /etc/resolv.conf
    ```

3. Install docker-compose tool by following this  
   https://docs.docker.com/compose/install/#install-compose