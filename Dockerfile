# 使用Busybox作为基础镜像，体积小但包含shell
FROM busybox:latest

# 创建输出目录
RUN mkdir -p /output

# 添加模拟脚本
COPY scripts/sublist3r-simulator.sh /app/sublist3r.sh
RUN chmod +x /app/sublist3r.sh

# 设置工作目录
WORKDIR /output

# 设置容器启动命令
ENTRYPOINT ["/app/sublist3r.sh"]
