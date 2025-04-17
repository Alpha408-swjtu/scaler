import matplotlib.pyplot as plt
import numpy as np
import os

# 日志文件路径
log_path = "./file"

# 手动指定文件名列表
file_names = ["real","double_avg_smoothing","single_exp_smoothing","triple_exp_smoothing"]  # 根据实际情况修改文件名

# 生成完整的文件路径列表
log_files = [os.path.join(log_path, f) for f in file_names]

# 用于存储每个文件的数据
data = {}

# 读取每个文件的内容
for log_file in log_files:
    with open(log_file, "r") as file:
        # 读取每行数据并转为 float
        data[os.path.basename(log_file)] = [float(line.strip()) for line in file.readlines()]

# 提取前360个数据点（3小时，每30秒一个数据点）
for file_name in data:
    data[file_name] = data[file_name][:400]

# 计算时间步数 (总共有360个数据点，代表3小时)
time_steps = np.arange(0, len(data[list(data.keys())[0]]))  # 每个时间步代表一个数据点

# 将横轴显示为0到360，代表3小时（每个数据点代表30秒）
time_labels = np.arange(0, len(time_steps), 30)  # 每30个数据点作为一个标签（对应15分钟）

# 绘制图形
plt.figure(figsize=(10, 6))

# 为每个文件绘制数据
for file_name in data:
    plt.plot(time_steps, data[file_name], label=f"{file_name}")

# 添加标题和标签
plt.title("Real and Predict Qps")
plt.xlabel("Timeline in 2.5 hours (Sampling every 30s)")
plt.ylabel("Value/")

# 设置横轴的刻度，每30个数据点为一个刻度
plt.xticks(time_labels, time_labels * 0.5)  # 将每30个数据点显示为分钟

# 旋转横轴标签以避免拥挤
plt.xticks(rotation=45)

# 添加图例
plt.legend(loc='upper right')

# 调整布局以避免标签被截断
plt.tight_layout()

# 保存图片
plt.savefig("load.png", dpi=300)

# 显示图像
plt.show()