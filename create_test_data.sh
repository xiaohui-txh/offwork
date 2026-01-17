#!/bin/bash
set -e

# 接口地址
API_URL="http://localhost:8080/api/v1/offwork/checkin"

# 中心点经纬度（例如北京市天安门）
CENTER_LNG=116.397
CENTER_LAT=39.908

# 打卡数量
COUNT=50

echo "开始生成 $COUNT 条测试数据 ..."

for i in $(seq 1 $COUNT); do
    # 生成随机偏移（经纬度 ±0.02 度，约 2km 范围）
    OFFSET_LNG=$(awk -v min=-0.02 -v max=0.02 'BEGIN{srand(); print min+rand()*(max-min)}')
    OFFSET_LAT=$(awk -v min=-0.02 -v max=0.02 'BEGIN{srand(); print min+rand()*(max-min)}')

    Lng=$(awk -v center=$CENTER_LNG -v offset=$OFFSET_LNG 'BEGIN{printf "%.6f", center+offset}')
    Lat=$(awk -v center=$CENTER_LAT -v offset=$OFFSET_LAT 'BEGIN{printf "%.6f", center+offset}')

    # 调用接口
    curl -s -X POST $API_URL \
        -H "Content-Type: application/json" \
        -d "{\"lng\":$Lng,\"lat\":$Lat}" > /dev/null

    echo "生成第 $i 条: lng=$Lng, lat=$Lat"
done

echo "测试数据生成完成"
