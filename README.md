# jsyslog

After years of using, managing, and building logging platforms,
with this experience, I wanted to build a system.

The [syslog protocol](https://en.wikipedia.org/wiki/Syslog) has been around, implemented, and supported in many different platforms.

## Usage

### Forwarding messages

* listen on TCP :9000 -> forward message to TCP 10.10.10.10:9000

  ```bash
  jsyslog forwarder \
    -from tcp://0.0.0.0:9000 \
    -to tcp://10.10.10.10:9000
  ```

* listen on UDP :9000 -> forward message to TCP 10.10.10.10:9000

  ```bash
  jsyslog forwarder \
    -from udp://0.0.0.0:9000 \
    -to tcp://10.10.10.10:9000
  ```

* listen on TCP & UDP :9000 -> write to local file

  ```bash
  jsyslog forwarder \
    -from tcp://0.0.0.0:9000 \
    -from udp://0.0.0.0:9000 \
    -to file:///var/log/syslog.log
  ```
