FROM ubuntu:22.04
RUN apt-get update && apt-get install -y stress
CMD ["stress", "--cpu", "4", "--vm", "2", "--vm-bytes", "256M", "--timeout", "300s"]
