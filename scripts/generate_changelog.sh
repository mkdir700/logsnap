#!/bin/bash

# 创建或清空 CHANGELOG.md 文件
echo "# Changelog" > CHANGELOG.md
echo "" >> CHANGELOG.md
echo "所有项目的显著变更都将记录在此文件中。" >> CHANGELOG.md
echo "" >> CHANGELOG.md
echo "格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，" >> CHANGELOG.md
echo "并且本项目遵循 [语义化版本](https://semver.org/lang/zh-CN/)。" >> CHANGELOG.md
echo "" >> CHANGELOG.md

# 获取所有标签并按版本排序
tags=$(git tag --sort=-v:refname)
tags_array=($tags)

# 如果没有标签，则使用最新提交
if [ ${#tags_array[@]} -eq 0 ]; then
    echo "## [未发布]" >> CHANGELOG.md
    echo "" >> CHANGELOG.md
    git log --pretty=format:"- %s (%h)" >> CHANGELOG.md
    exit 0
fi

# 处理每个标签
for ((i=0; i<${#tags_array[@]}; i++)); do
    current_tag=${tags_array[$i]}
    
    # 获取标签日期
    tag_date=$(git log -1 --format=%ad --date=short $current_tag)
    
    echo "## [$current_tag] - $tag_date" >> CHANGELOG.md
    echo "" >> CHANGELOG.md
    
    # 确定比较的范围
    if [ $i -eq $((${#tags_array[@]}-1)) ]; then
        # 最后一个标签，从项目开始到此标签
        git log $current_tag --pretty=format:"- %s (%h)" --no-merges | grep -v "Merge" >> CHANGELOG.md
    else
        # 从上一个标签到当前标签
        next_tag=${tags_array[$i+1]}
        git log $next_tag..$current_tag --pretty=format:"- %s (%h)" --no-merges | grep -v "Merge" >> CHANGELOG.md
    fi
    
    echo "" >> CHANGELOG.md
    echo "" >> CHANGELOG.md
done

echo "Changelog 已生成到 CHANGELOG.md 文件中" 