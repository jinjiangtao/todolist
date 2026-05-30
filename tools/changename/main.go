package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
)

type FileInfo struct {
	Path    string
	Name    string
	Ext     string
	NewName string
}

type Rule struct {
	Type     string // "sequence", "replace", "date", "case", "extension"
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
		Source    string // "modified"
		Format    string // "YYYYMMDD", "YYYY-MM-DD" etc
		Position  string // "prefix", "suffix"
		Separator string
	}
	Case struct {
		Mode string // "upper", "lower", "title", "camel"
	}
	Extension struct {
		NewExt    string
		UnifyCase bool
		TargetCase string
	}
}

type AppState struct {
	files       []FileInfo
	rules       []Rule
	undoRecords map[string][]UndoItem
	mu          sync.Mutex
}

type UndoItem struct {
	OldPath string
	NewPath string
}

func main() {
	a := app.New()
	a.Settings().SetTheme(theme.DefaultTheme())
	w := a.NewWindow("智能文件批量重命名工具 - Smart Renamer")
	w.Resize(fyne.NewSize(1200, 800))

	state := &AppState{
		files:       []FileInfo{},
		rules:       createDefaultRules(),
		undoRecords: make(map[string][]UndoItem),
	}

	ui := createUI(w, state)
	w.SetContent(ui)
	w.ShowAndRun()
}

func createDefaultRules() []Rule {
	return []Rule{
		{Type: "sequence", Enabled: false},
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

func createUI(w fyne.Window, state *AppState) fyne.CanvasObject {
	// 左侧面板 - 规则配置
	rulePanel := createRulePanel(state, w)
	
	// 右侧面板 - 文件列表和操作
	filePanel := createFilePanel(state, w)
	
	// 主布局
	split := container.NewHSplit(rulePanel, filePanel)
	split.Offset = 0.35
	
	return split
}

func createRulePanel(state *AppState, w fyne.Window) fyne.CanvasObject {
	content := container.NewVBox()
	
	header := widget.NewLabelWithStyle("重命名规则", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	content.Add(header)
	content.Add(widget.NewSeparator())
	
	ruleTypes := map[string]string{
		"sequence":    "序号填充",
		"replace":     "查找替换",
		"date":        "日期标记",
		"case":        "大小写转换",
		"extension":   "扩展名修改",
	}
	
	for i := range state.rules {
		rule := &state.rules[i]
		name := ruleTypes[rule.Type]
		
		ruleCard := createRuleCard(rule, name, state, w)
		content.Add(ruleCard)
		content.Add(widget.NewSeparator())
	}
	
	return container.NewScroll(content)
}

func createRuleCard(rule *Rule, name string, state *AppState, w fyne.Window) fyne.CanvasObject {
	check := widget.NewCheck(name, func(checked bool) {
		rule.Enabled = checked
		updatePreview(state)
	})
	check.SetChecked(rule.Enabled)
	
	content := container.NewVBox()
	
	switch rule.Type {
	case "sequence":
		startEntry := widget.NewEntry()
		startEntry.SetText(fmt.Sprintf("%d", rule.Sequence.Start))
		startEntry.SetPlaceHolder("起始数字")
		startEntry.OnChanged = func(s string) {
			if val, err := strconv.Atoi(s); err == nil {
				rule.Sequence.Start = val
				updatePreview(state)
			}
		}
		
		digitsEntry := widget.NewEntry()
		digitsEntry.SetText(fmt.Sprintf("%d", rule.Sequence.Digits))
		digitsEntry.SetPlaceHolder("位数")
		digitsEntry.OnChanged = func(s string) {
			if val, err := strconv.Atoi(s); err == nil {
				rule.Sequence.Digits = val
				updatePreview(state)
			}
		}
		
		prefixEntry := widget.NewEntry()
		prefixEntry.SetText(rule.Sequence.Prefix)
		prefixEntry.SetPlaceHolder("前缀")
		prefixEntry.OnChanged = func(s string) {
			rule.Sequence.Prefix = s
			updatePreview(state)
		}
		
		suffixEntry := widget.NewEntry()
		suffixEntry.SetText(rule.Sequence.Suffix)
		suffixEntry.SetPlaceHolder("后缀")
		suffixEntry.OnChanged = func(s string) {
			rule.Sequence.Suffix = s
			updatePreview(state)
		}
		
		content.Add(widget.NewLabel("起始数字:"))
		content.Add(startEntry)
		content.Add(widget.NewLabel("位数:"))
		content.Add(digitsEntry)
		content.Add(widget.NewLabel("前缀:"))
		content.Add(prefixEntry)
		content.Add(widget.NewLabel("后缀:"))
		content.Add(suffixEntry)
		
	case "replace":
		searchEntry := widget.NewEntry()
		searchEntry.SetText(rule.Replace.Search)
		searchEntry.SetPlaceHolder("查找内容")
		searchEntry.OnChanged = func(s string) {
			rule.Replace.Search = s
			updatePreview(state)
		}
		
		replaceEntry := widget.NewEntry()
		replaceEntry.SetText(rule.Replace.Replace)
		replaceEntry.SetPlaceHolder("替换为")
		replaceEntry.OnChanged = func(s string) {
			rule.Replace.Replace = s
			updatePreview(state)
		}
		
		regexCheck := widget.NewCheck("使用正则表达式", func(checked bool) {
			rule.Replace.UseRegex = checked
			updatePreview(state)
		})
		regexCheck.SetChecked(rule.Replace.UseRegex)
		
		content.Add(widget.NewLabel("查找内容:"))
		content.Add(searchEntry)
		content.Add(widget.NewLabel("替换为:"))
		content.Add(replaceEntry)
		content.Add(regexCheck)
		
	case "date":
		formatSelect := widget.NewSelect([]string{"YYYYMMDD", "YYYY-MM-DD", "YYYY_MM_DD", "YYYYMMDD_HHMMSS"}, func(s string) {
			rule.Date.Format = s
			updatePreview(state)
		})
		formatSelect.SetSelected(rule.Date.Format)
		
		positionSelect := widget.NewSelect([]string{"prefix", "suffix"}, func(s string) {
			rule.Date.Position = s
			updatePreview(state)
		})
		positionSelect.SetSelected(rule.Date.Position)
		
		sepEntry := widget.NewEntry()
		sepEntry.SetText(rule.Date.Separator)
		sepEntry.SetPlaceHolder("分隔符")
		sepEntry.OnChanged = func(s string) {
			rule.Date.Separator = s
			updatePreview(state)
		}
		
		content.Add(widget.NewLabel("日期格式:"))
		content.Add(formatSelect)
		content.Add(widget.NewLabel("位置:"))
		content.Add(positionSelect)
		content.Add(widget.NewLabel("分隔符:"))
		content.Add(sepEntry)
		
	case "case":
		caseSelect := widget.NewSelect([]string{"upper", "lower", "title", "camel"}, func(s string) {
			rule.Case.Mode = s
			updatePreview(state)
		})
		caseSelect.SetSelected(rule.Case.Mode)
		
		content.Add(widget.NewLabel("转换方式:"))
		content.Add(caseSelect)
		
	case "extension":
		extEntry := widget.NewEntry()
		extEntry.SetText(rule.Extension.NewExt)
		extEntry.SetPlaceHolder("新扩展名 (留空不修改)")
		extEntry.OnChanged = func(s string) {
			rule.Extension.NewExt = s
			updatePreview(state)
		}
		
		unifyCheck := widget.NewCheck("统一扩展名大小写", func(checked bool) {
			rule.Extension.UnifyCase = checked
			updatePreview(state)
		})
		unifyCheck.SetChecked(rule.Extension.UnifyCase)
		
		targetCaseSelect := widget.NewSelect([]string{"lower", "upper"}, func(s string) {
			rule.Extension.TargetCase = s
			updatePreview(state)
		})
		targetCaseSelect.SetSelected("lower")
		
		content.Add(widget.NewLabel("新扩展名:"))
		content.Add(extEntry)
		content.Add(unifyCheck)
		content.Add(widget.NewLabel("目标格式:"))
		content.Add(targetCaseSelect)
	}
	
	form := container.NewVBox(check)
	ruleContent := container.NewVBox(content...)
	ruleContent.Hide()
	
	check.OnChanged = func(checked bool) {
		rule.Enabled = checked
		if checked {
			ruleContent.Show()
		} else {
			ruleContent.Hide()
		}
		updatePreview(state)
	}
	
	return container.NewVBox(form, ruleContent)
}

func createFilePanel(state *AppState, w fyne.Window) fyne.CanvasObject {
	// 顶部工具栏
	toolbar := container.NewBorder(nil, nil, nil, nil, container.NewHBox(
		widget.NewButtonWithIcon("选择文件夹", theme.FolderOpenIcon(), func() {
			openFolderDialog(state, w)
		}),
		widget.NewButton("刷新预览", func() {
			updatePreview(state)
		}),
	))
	
	// 文件列表
	list := widget.NewTable(
		func() (int, int) {
			return len(state.files), 3
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			if id.Row < len(state.files) {
				file := state.files[id.Row]
				switch id.Col {
				case 0:
					label.SetText(file.Name + file.Ext)
				case 1:
					label.SetText("→")
				case 2:
					if file.NewName != "" && file.NewName != file.Name+file.Ext {
						label.SetText(file.NewName)
					} else {
						label.SetText("-")
					}
				}
			}
		},
	)
	list.SetColumnWidth(0, 300)
	list.SetColumnWidth(1, 30)
	list.SetColumnWidth(2, 300)
	
	// 底部操作栏
	statusLabel := widget.NewLabel("就绪")
	
	undoBtn := widget.NewButtonWithIcon("撤销", theme.MediaReplayIcon(), func() {
		showUndoDialog(state, w)
	})
	undoBtn.Disable()
	
	executeBtn := widget.NewButtonWithIcon("开始重命名", theme.ConfirmIcon(), func() {
		executeRename(state, w, statusLabel, undoBtn)
	})
	executeBtn.Importance = widget.HighImportance
	
	bottom := container.NewBorder(nil, nil, nil, container.NewHBox(undoBtn, executeBtn), statusLabel)
	
	return container.NewBorder(toolbar, bottom, nil, nil, list)
}

func openFolderDialog(state *AppState, w fyne.Window) {
	dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
		if err != nil || uri == nil {
			return
		}
		
		path := uri.Path()
		scanFiles(path, state, w)
	}, w)
}

func scanFiles(path string, state *AppState, w fyne.Window) {
	var files []FileInfo
	var mu sync.Mutex
	var wg sync.WaitGroup

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
		dialog.ShowError(err, w)
		return
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name+files[i].Ext < files[j].Name+files[j].Ext
	})

	state.files = files
	updatePreview(state)
}

func updatePreview(state *AppState) {
	for i := range state.files {
		state.files[i].NewName = applyRules(state.files[i], state.rules, i)
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

func executeRename(state *AppState, w fyne.Window, status *widget.Label, undoBtn *widget.Button) {
	if len(state.files) == 0 {
		dialog.ShowInformation("提示", "请先选择文件夹", w)
		return
	}

	confirm := dialog.NewConfirm("确认重命名", "确定要重命名这些文件吗？", func(ok bool) {
		if !ok {
			return
		}

		var undoItems []UndoItem
		renamed := 0
		skipped := 0

		processedPaths := make(map[string]bool)

		for _, file := range state.files {
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
			state.undoRecords[undoID] = undoItems
			undoBtn.Enable()
		}

		status.SetText(fmt.Sprintf("已重命名 %d 个文件，跳过 %d 个", renamed, skipped))
		
		// 重新扫描
		if len(state.files) > 0 {
			scanFiles(filepath.Dir(state.files[0].Path), state, w)
		}
	}, w)
	confirm.Show()
}

func showUndoDialog(state *AppState, w fyne.Window) {
	var undoID string
	var items []UndoItem
	
	for id, u := range state.undoRecords {
		undoID = id
		items = u
		break
	}
	
	if len(items) == 0 {
		dialog.ShowInformation("提示", "没有可撤销的操作", w)
		return
	}

	confirm := dialog.NewConfirm("撤销操作", fmt.Sprintf("确定要撤销上次对 %d 个文件的重命名操作吗？", len(items)), func(ok bool) {
		if !ok {
			return
		}
		
		restored := 0
		for _, item := range items {
			if err := os.Rename(item.OldPath, item.NewPath); err == nil {
				restored++
			}
		}
		
		delete(state.undoRecords, undoID)
		dialog.ShowInformation("完成", fmt.Sprintf("成功恢复 %d 个文件", restored), w)
		
		// 重新扫描
		if len(state.files) > 0 {
			scanFiles(filepath.Dir(state.files[0].Path), state, w)
		}
	}, w)
	confirm.Show()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
