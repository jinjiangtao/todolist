package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type FileInfo struct {
	Path    string
	Name    string
	Ext     string
	NewName string
}

type Rule struct {
	Type     string
	Enabled  bool
	Sequence struct {
		Start  int
		Digits int
		Prefix string
		Suffix string
	}
	Replace struct {
		Search   string
		Replace  string
		UseRegex bool
	}
	Date struct {
		Source    string
		Format    string
		Position  string
		Separator string
	}
	Case struct {
		Mode string
	}
	Extension struct {
		NewExt     string
		UnifyCase  bool
		TargetCase string
	}
}

type UndoItem struct {
	OldPath string
	NewPath string
}

var (
	files       []FileInfo
	rules       []Rule
	undoRecords map[string][]UndoItem
	reader      = bufio.NewReader(os.Stdin)
)

func main() {
	rules = createDefaultRules()
	undoRecords = make(map[string][]UndoItem)

	println("========================================")
	println("  智能文件批量重命名工具 - Smart Renamer")
	println("========================================")

	for {
		println("\n请选择操作:")
		println("1. 选择文件夹")
		println("2. 配置重命名规则")
		println("3. 预览重命名结果")
		println("4. 执行重命名")
		println("5. 撤销上次操作")
		println("6. 退出")
		print("请输入选项 (1-6): ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			selectFolder()
		case "2":
			configRules()
		case "3":
			previewRename()
		case "4":
			executeRename()
		case "5":
			undoRename()
		case "6":
			println("再见！")
			return
		default:
			println("无效选项，请重试。")
		}
	}
}

func createDefaultRules() []Rule {
	return []Rule{
		{Type: "sequence", Enabled: false, Sequence: struct {
			Start  int
			Digits int
			Prefix string
			Suffix string
		}{Start: 1, Digits: 3}},
		{Type: "replace", Enabled: false},
		{Type: "date", Enabled: false, Date: struct {
			Source    string
			Format    string
			Position  string
			Separator string
		}{Source: "modified", Format: "YYYYMMDD", Position: "prefix", Separator: "_"}},
		{Type: "case", Enabled: false, Case: struct{ Mode string }{Mode: "lower"}},
		{Type: "extension", Enabled: false},
	}
}

func selectFolder() {
	print("\n请输入文件夹路径: ")
	path, _ := reader.ReadString('\n')
	path = strings.TrimSpace(path)
	path = strings.Trim(path, "\"")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		println("错误: 文件夹不存在")
		return
	}

	scanFiles(path)
	println(fmt.Sprintf("成功加载 %d 个文件", len(files)))
}

func scanFiles(path string) {
	var mu sync.Mutex
	var wg sync.WaitGroup

	files = nil

	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if p == path {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			ext := filepath.Ext(p)
			name := info.Name()
			if ext != "" {
				name = name[:len(name)-len(ext)]
			}
			mu.Lock()
			files = append(files, FileInfo{
				Path: p,
				Name: name,
				Ext:  ext,
			})
			mu.Unlock()
		}()

		return nil
	})

	wg.Wait()

	if err != nil {
		println("扫描文件夹时出错:", err.Error())
		return
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name+files[i].Ext < files[j].Name+files[j].Ext
	})
}

func configRules() {
	for {
		println("\n当前规则状态:")
		for i, rule := range rules {
			status := "关闭"
			if rule.Enabled {
				status = "开启"
			}
			ruleName := getRuleName(rule.Type)
			println(fmt.Sprintf("%d. %s [%s]", i+1, ruleName, status))
		}
		println("0. 返回主菜单")
		print("请选择要配置的规则 (0-5): ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		
		if input == "0" {
			return
		}

		idx, err := strconv.Atoi(input)
		if err != nil || idx < 1 || idx > 5 {
			println("无效选项")
			continue
		}

		configRule(idx - 1)
	}
}

func getRuleName(t string) string {
	switch t {
	case "sequence":
		return "序号填充"
	case "replace":
		return "查找替换"
	case "date":
		return "日期标记"
	case "case":
		return "大小写转换"
	case "extension":
		return "扩展名修改"
	default:
		return t
	}
}

func configRule(idx int) {
	rule := &rules[idx]
	
	println(fmt.Sprintf("\n配置规则: %s", getRuleName(rule.Type)))
	print("是否启用该规则? (y/n): ")
	ans, _ := reader.ReadString('\n')
	rule.Enabled = strings.ToLower(strings.TrimSpace(ans)) == "y"

	if !rule.Enabled {
		return
	}

	switch rule.Type {
	case "sequence":
		print("起始数字 (默认 1): ")
		input, _ := reader.ReadString('\n')
		if v, err := strconv.Atoi(strings.TrimSpace(input)); err == nil {
			rule.Sequence.Start = v
		}

		print("位数 (默认 3): ")
		input, _ = reader.ReadString('\n')
		if v, err := strconv.Atoi(strings.TrimSpace(input)); err == nil {
			rule.Sequence.Digits = v
		}

		print("前缀 (默认空): ")
		input, _ = reader.ReadString('\n')
		rule.Sequence.Prefix = strings.TrimSpace(input)

		print("后缀 (默认空): ")
		input, _ = reader.ReadString('\n')
		rule.Sequence.Suffix = strings.TrimSpace(input)

	case "replace":
		print("查找内容: ")
		input, _ := reader.ReadString('\n')
		rule.Replace.Search = strings.TrimSpace(input)

		print("替换为: ")
		input, _ = reader.ReadString('\n')
		rule.Replace.Replace = strings.TrimSpace(input)

		print("使用正则表达式? (y/n): ")
		input, _ = reader.ReadString('\n')
		rule.Replace.UseRegex = strings.ToLower(strings.TrimSpace(input)) == "y"

	case "date":
		println("日期格式:")
		println("1. YYYYMMDD")
		println("2. YYYY-MM-DD")
		println("3. YYYY_MM_DD")
		println("4. YYYYMMDD_HHMMSS")
		print("请选择格式 (1-4): ")
		input, _ := reader.ReadString('\n')
		switch strings.TrimSpace(input) {
		case "1":
			rule.Date.Format = "YYYYMMDD"
		case "2":
			rule.Date.Format = "YYYY-MM-DD"
		case "3":
			rule.Date.Format = "YYYY_MM_DD"
		case "4":
			rule.Date.Format = "YYYYMMDD_HHMMSS"
		}

		print("位置 (prefix/suffix, 默认 prefix): ")
		input, _ = reader.ReadString('\n')
		rule.Date.Position = strings.TrimSpace(input)
		if rule.Date.Position != "prefix" && rule.Date.Position != "suffix" {
			rule.Date.Position = "prefix"
		}

		print("分隔符 (默认 _): ")
		input, _ = reader.ReadString('\n')
		rule.Date.Separator = strings.TrimSpace(input)
		if rule.Date.Separator == "" {
			rule.Date.Separator = "_"
		}

	case "case":
		println("转换方式:")
		println("1. upper (全大写)")
		println("2. lower (全小写)")
		println("3. title (首字母大写)")
		println("4. camel (驼峰式)")
		print("请选择 (1-4): ")
		input, _ := reader.ReadString('\n')
		switch strings.TrimSpace(input) {
		case "1":
			rule.Case.Mode = "upper"
		case "2":
			rule.Case.Mode = "lower"
		case "3":
			rule.Case.Mode = "title"
		case "4":
			rule.Case.Mode = "camel"
		}

	case "extension":
		print("新扩展名 (留空不修改): ")
		input, _ := reader.ReadString('\n')
		rule.Extension.NewExt = strings.TrimSpace(input)

		print("统一扩展名大小写? (y/n): ")
		input, _ = reader.ReadString('\n')
		rule.Extension.UnifyCase = strings.ToLower(strings.TrimSpace(input)) == "y"

		if rule.Extension.UnifyCase {
			print("目标格式 (lower/upper, 默认 lower): ")
			input, _ = reader.ReadString('\n')
			rule.Extension.TargetCase = strings.TrimSpace(input)
			if rule.Extension.TargetCase != "lower" && rule.Extension.TargetCase != "upper" {
				rule.Extension.TargetCase = "lower"
			}
		}
	}

	println("规则配置完成")
}

func previewRename() {
	if len(files) == 0 {
		println("请先选择文件夹")
		return
	}

	enabledCount := 0
	for _, rule := range rules {
		if rule.Enabled {
			enabledCount++
		}
	}

	if enabledCount == 0 {
		println("请先配置并启用至少一个重命名规则")
		return
	}

	for i := range files {
		files[i].NewName = applyRules(files[i], rules, i)
	}

	println("\n预览结果:")
	println("----------------------------------------------------------------------")
	println(fmt.Sprintf("%-30s %-30s %-10s", "原文件名", "新文件名", "状态"))
	println("----------------------------------------------------------------------")

	for _, file := range files {
		status := "不变"
		if file.NewName != file.Name+file.Ext {
			status = "重命名"
		}
		println(fmt.Sprintf("%-30s %-30s %-10s", 
			truncate(file.Name+file.Ext, 28), 
			truncate(file.NewName, 28), 
			status))
	}

	println("----------------------------------------------------------------------")
	println(fmt.Sprintf("共 %d 个文件，其中 %d 个将被重命名", 
		len(files), countToRename()))
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func countToRename() int {
	count := 0
	for _, file := range files {
		if file.NewName != "" && file.NewName != file.Name+file.Ext {
			count++
		}
	}
	return count
}

func executeRename() {
	if len(files) == 0 {
		println("请先选择文件夹")
		return
	}

	if countToRename() == 0 {
		println("没有文件需要重命名")
		return
	}

	print("\n确定要执行重命名吗? (y/n): ")
	ans, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(ans)) != "y" {
		println("操作取消")
		return
	}

	var undoItems []UndoItem
	renamed := 0
	skipped := 0

	processedPaths := make(map[string]bool)

	for _, file := range files {
		oldPath := file.Path
		newName := file.NewName
		if newName == "" || newName == file.Name+file.Ext {
			continue
		}

		newPath := filepath.Join(filepath.Dir(oldPath), newName)
		
		finalPath := newPath
		counter := 1
		for (fileExists(finalPath) && oldPath != finalPath) || processedPaths[finalPath] {
			ext := filepath.Ext(finalPath)
			base := finalPath[:len(finalPath)-len(ext)]
			finalPath = fmt.Sprintf("%s_%d%s", base, counter, ext)
			counter++
		}

		if err := os.Rename(oldPath, finalPath); err != nil {
			println(fmt.Sprintf("重命名失败: %s -> %s, 错误: %s", 
				file.Name+file.Ext, newName, err.Error()))
			skipped++
			continue
		}

		undoItems = append(undoItems, UndoItem{
			OldPath: finalPath,
			NewPath: oldPath,
		})
		processedPaths[finalPath] = true
		renamed++
	}

	if len(undoItems) > 0 {
		undoID := uuid.New().String()
		undoRecords[undoID] = undoItems
	}

	println(fmt.Sprintf("完成! 重命名 %d 个文件，跳过 %d 个", renamed, skipped))
	
	if len(files) > 0 {
		scanFiles(filepath.Dir(files[0].Path))
	}
}

func undoRename() {
	var undoID string
	var items []UndoItem
	
	for id, u := range undoRecords {
		undoID = id
		items = u
		break
	}
	
	if len(items) == 0 {
		println("没有可撤销的操作")
		return
	}

	print(fmt.Sprintf("确定要撤销上次对 %d 个文件的重命名操作吗? (y/n): ", len(items)))
	ans, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(ans)) != "y" {
		println("操作取消")
		return
	}

	restored := 0
	for _, item := range items {
		if err := os.Rename(item.OldPath, item.NewPath); err == nil {
			restored++
		}
	}

	delete(undoRecords, undoID)
	println(fmt.Sprintf("完成! 成功恢复 %d 个文件", restored))
	
	if len(files) > 0 {
		scanFiles(filepath.Dir(files[0].Path))
	}
}

func applyRules(file FileInfo, rules []Rule, index int) string {
	name := file.Name
	ext := file.Ext

	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}

		switch rule.Type {
		case "sequence":
			num := rule.Sequence.Start + index
			format := fmt.Sprintf("%%0%dd", rule.Sequence.Digits)
			name = rule.Sequence.Prefix + fmt.Sprintf(format, num) + rule.Sequence.Suffix
		case "replace":
			if rule.Replace.UseRegex {
				re, err := regexp.Compile(rule.Replace.Search)
				if err == nil {
					name = re.ReplaceAllString(name, rule.Replace.Replace)
				}
			} else {
				name = strings.ReplaceAll(name, rule.Replace.Search, rule.Replace.Replace)
			}
		case "date":
			dateStr := getDateString(file.Path, rule.Date.Format)
			if rule.Date.Position == "prefix" {
				name = dateStr + rule.Date.Separator + name
			} else {
				name = name + rule.Date.Separator + dateStr
			}
		case "case":
			name = applyCaseConversion(name, rule.Case.Mode)
		case "extension":
			ext = applyExtensionChange(ext, rule.Extension)
		}
	}

	if ext != "" {
		return name + ext
	}
	return name
}

func getDateString(path, format string) string {
	info, err := os.Stat(path)
	if err != nil {
		return time.Now().Format("20060102")
	}
	t := info.ModTime()
	
	switch format {
	case "YYYYMMDD":
		return t.Format("20060102")
	case "YYYY-MM-DD":
		return t.Format("2006-01-02")
	case "YYYY_MM_DD":
		return t.Format("2006_01_02")
	case "YYYYMMDD_HHMMSS":
		return t.Format("20060102_150405")
	default:
		return t.Format("20060102")
	}
}

func applyCaseConversion(s, mode string) string {
	switch mode {
	case "upper":
		return strings.ToUpper(s)
	case "lower":
		return strings.ToLower(s)
	case "title":
		if len(s) == 0 {
			return s
		}
		return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
	case "camel":
		words := strings.FieldsFunc(s, func(r rune) bool {
			return r == ' ' || r == '_' || r == '-'
		})
		if len(words) == 0 {
			return s
		}
		result := strings.ToLower(words[0])
		for i := 1; i < len(words); i++ {
			result += strings.Title(words[i])
		}
		return result
	default:
		return s
	}
}

func applyExtensionChange(ext string, rule struct {
	NewExt     string
	UnifyCase  bool
	TargetCase string
}) string {
	if rule.NewExt != "" {
		newExt := rule.NewExt
		if !strings.HasPrefix(newExt, ".") {
			newExt = "." + newExt
		}
		ext = newExt
	}

	if rule.UnifyCase {
		if rule.TargetCase == "lower" {
			ext = strings.ToLower(ext)
		} else {
			ext = strings.ToUpper(ext)
		}
	}

	return ext
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
