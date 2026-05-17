#!/bin/bash
# GoAWK 功能测试脚本

echo "========================================="
echo "GoAWK 功能测试"
echo "========================================="
echo ""

# 测试1: 打印第一列
echo "测试1: 打印第一列"
echo "命令: awk '\$1' test.txt"
echo "结果:"
awk '$1' test.txt
echo ""

# 测试2: 条件过滤 - 数值比较
echo "测试2: 条件过滤 - 数值比较"
echo "命令: awk '\$2 > 28 { print \$1 }' test.txt"
echo "结果:"
awk '$2 > 28 { print $1 }' test.txt
echo ""

# 测试3: 字符串比较
echo "测试3: 字符串比较"
echo "命令: awk '\$1 == \"Alice\" { print \$2, \$3 }' test.txt"
echo "结果:"
awk '$1 == "Alice" { print $2, $3 }' test.txt
echo ""

# 测试4: 管道输入
echo "测试4: 管道输入"
echo "命令: echo \"hello world go\" | awk '\$1'"
echo "结果:"
echo "hello world go" | awk '$1'
echo ""

# 测试5: CSV 处理
echo "测试5: CSV 处理"
echo "命令: awk -F ',' '\$2' test.csv"
echo "结果:"
awk -F ',' '$2' test.csv
echo ""

# 测试6: 多列组合
echo "测试6: 多列组合"
echo "命令: awk '{ print \$1, \"-\", \$3 }' test.txt"
echo "结果:"
awk '{ print $1, "-", $3 }' test.txt
echo ""

# 测试7: 中文处理
echo "测试7: 中文处理"
echo "命令: awk '\$1' test_cn.txt"
echo "结果:"
awk '$1' test_cn.txt
echo ""

# 测试8: 跳过表头
echo "测试8: 跳过表头"
echo "命令: awk 'NR>1 { print \$1, \$2 }' test.csv"
echo "结果:"
awk 'NR>1 { print $1, $2 }' test.csv
echo ""

echo "========================================="
echo "所有测试完成"
echo "========================================="
