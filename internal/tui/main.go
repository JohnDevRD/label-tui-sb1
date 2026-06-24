package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/JohnDevRD/label-tui-sb1/internal/core"
	"github.com/JohnDevRD/label-tui-sb1/internal/printers"
)

type screen int

const (
	screenWelcome screen = iota
	screenLogin
	screenSearch
	screenSelectTemplate
	screenPreview
	screenSettings
)

type Model struct {
	current     screen
	sapClient   *core.SapClient
	settings    *core.Settings
	articles    []core.Article
	printJobs   []core.PrintJob
	templates   []string
	width       int
	height      int

	welcome     welcomeModel
	login       loginModel
	search      searchModel
	tmplSelect  tmplSelectModel
	preview     previewModel
	settingsScr settingsModel
}

func New() Model {
	settings, _ := core.LoadSettings()
	if settings == nil {
		settings = &core.Settings{}
	}

	m := Model{
		current:  screenWelcome,
		settings: settings,
		welcome:  newWelcomeModel(),
		login:    newLoginModel(),
	}

	if settings.SAPServiceLayerURL != "" {
		m.sapClient = core.NewSapClient(settings.SAPServiceLayerURL, settings.PriceList)
	}

	return m
}

func (m *Model) SetTemplates(templates []string) {
	m.templates = templates
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.welcome.width = msg.Width
		m.welcome.height = msg.Height
		m.search.width = msg.Width
		m.search.height = msg.Height
		m.settingsScr.width = msg.Width
		m.settingsScr.height = msg.Height
		return m, nil

	case loginResultMsg:
		if msg.success {
			m.current = screenSearch
			m.search = newSearchModel()
			m.tmplSelect = newTmplSelectModel()
			return m, nil
		}
		m.login.err = msg.err
		return m, nil

	case settingsSavedMsg:
		if msg.err != "" {
			m.settingsScr.err = msg.err
		} else {
			m.settingsScr.saved = true
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if m.current == screenWelcome {
				return m, tea.Quit
			}
			m.current = screenWelcome
			return m, nil
		}
	}

	switch m.current {
	case screenWelcome:
		return m.updateWelcome(msg)
	case screenLogin:
		return m.updateLogin(msg)
	case screenSearch:
		return m.updateSearch(msg)
	case screenSelectTemplate:
		return m.updateTemplateSelect(msg)
	case screenPreview:
		return m.updatePreview(msg)
	case screenSettings:
		return m.updateSettings(msg)
	}

	return m, nil
}

func (m Model) View() string {
	w := m.width
	h := m.height
	if w == 0 {
		w = 80
	}
	if h == 0 {
		h = 24
	}

	header := TitleStyle.Render("Label TUI  •  SAP B1 Label Printer")

	var body string
	switch m.current {
	case screenWelcome:
		body = m.viewWelcome()
	case screenLogin:
		body = m.viewLogin()
	case screenSearch:
		body = m.viewSearch()
	case screenSelectTemplate:
		body = m.viewTemplateSelect()
	case screenPreview:
		body = m.viewPreview()
	case screenSettings:
		body = m.viewSettings()
	}

	content := lipgloss.JoinVertical(lipgloss.Center, header, body)

	return lipgloss.Place(w, h,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

type welcomeModel struct {
	width, height int
	selected      int
	choices       []string
}

func newWelcomeModel() welcomeModel {
	return welcomeModel{
		choices: []string{"Login to SAP B1", "Settings", "Exit"},
	}
}

func (wm *welcomeModel) update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if wm.selected > 0 {
				wm.selected--
			}
		case "down", "j":
			if wm.selected < len(wm.choices)-1 {
				wm.selected++
			}
		case "enter":
			return func() tea.Msg {
				return welcomeSelectedMsg{index: wm.selected}
			}
		}
	}
	return nil
}

type welcomeSelectedMsg struct {
	index int
}

func (m *Model) updateWelcome(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case welcomeSelectedMsg:
		switch msg.index {
		case 0:
			m.current = screenLogin
		case 1:
			m.current = screenSettings
			m.settingsScr = newSettingsModel(m.settings)
		case 2:
			return m, tea.Quit
		}
	}

	cmd := m.welcome.update(msg)
	return m, cmd
}

func (m *Model) viewWelcome() string {
	subtitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorMuted)).
		Render("Connect to SAP B1 and print labels to Zebra printers")

	menu := RenderMenu(m.welcome.choices, m.welcome.selected)

	return lipgloss.JoinVertical(lipgloss.Center,
		subtitle,
		"",
		menu,
		"",
		RenderHelp("↑/↓ navigate", "enter select", "esc quit"),
	)
}

type loginModel struct {
	inputs  []textinput.Model
	focused int
	err     string
	loading bool
}

func newLoginModel() loginModel {
	inputs := make([]textinput.Model, 2)

	inputs[0] = textinput.New()
	inputs[0].Placeholder = "Username"
	inputs[0].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorPrimary))
	inputs[0].Focus()

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "Password"
	inputs[1].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorPrimary))
	inputs[1].EchoMode = textinput.EchoPassword

	return loginModel{inputs: inputs}
}

func (m *Model) updateLogin(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.login.loading {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			m.login.focused = (m.login.focused + 1) % len(m.login.inputs)
			return m, m.focusLoginInput()
		case "shift+tab", "up":
			m.login.focused = (m.login.focused - 1 + len(m.login.inputs)) % len(m.login.inputs)
			return m, m.focusLoginInput()
		case "enter":
			return m, m.doLogin()
		}
	}

	cmd := m.updateLoginInputs(msg)
	return m, cmd
}

func (m *Model) focusLoginInput() tea.Cmd {
	var cmds []tea.Cmd
	for i := range m.login.inputs {
		if i == m.login.focused {
			cmds = append(cmds, m.login.inputs[i].Focus())
		} else {
			m.login.inputs[i].Blur()
		}
	}
	return tea.Batch(cmds...)
}

func (m *Model) updateLoginInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.login.inputs))
	for i := range m.login.inputs {
		m.login.inputs[i], cmds[i] = m.login.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m *Model) doLogin() tea.Cmd {
	return func() tea.Msg {
		user := m.login.inputs[0].Value()
		password := m.login.inputs[1].Value()

		if user == "" || password == "" {
			return loginResultMsg{err: "Username and password are required"}
		}

		if m.settings.CompanyDB == "" {
			return loginResultMsg{err: "CompanyDB not configured. Go to Settings first."}
		}

		if m.sapClient == nil {
			return loginResultMsg{err: "SAP Service Layer URL not configured. Go to Settings first."}
		}

		if err := m.sapClient.Login(m.settings.CompanyDB, user, password); err != nil {
			return loginResultMsg{err: err.Error()}
		}

		return loginResultMsg{success: true}
	}
}

type loginResultMsg struct {
	success bool
	err     string
}

func (m *Model) viewLogin() string {
	var inputs []string
	for i, input := range m.login.inputs {
		style := InputStyle
		if i == m.login.focused {
			style = InputFocusedStyle
		}
		inputs = append(inputs, style.Render(input.View()))
	}

	form := lipgloss.JoinVertical(lipgloss.Left, inputs...)

	var status string
	if m.login.loading {
		status = SpinnerStyle.Render("Logging in to SAP B1...")
	} else if m.login.err != "" {
		status = ErrorStyle.Render("Error: " + m.login.err)
	}

	return lipgloss.JoinVertical(lipgloss.Center,
		SectionStyle.Render("SAP B1 Credentials"),
		"",
		form,
		"",
		status,
		"",
		RenderHelp("tab/shift+tab navigate", "enter login", "esc back"),
	)
}

type searchModel struct {
	width, height int
	input         textinput.Model
	results       []core.Article
	cursor        int
	err           string
	loading       bool
	qty           int
	selected      map[int]int
}

func newSearchModel() searchModel {
	ti := textinput.New()
	ti.Placeholder = "Search articles by code or description..."
	ti.Prompt = "🔍 "
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorAccent))
	ti.Focus()
	return searchModel{
		input:    ti,
		selected: make(map[int]int),
		qty:      1,
	}
}

func (m *Model) updateSearch(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.search.loading {
		return m, nil
	}

	switch msg := msg.(type) {
	case searchResultsMsg:
		m.search.results = msg.results
		m.search.err = msg.err
		m.search.loading = false
		m.search.cursor = 0
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.search.input.Focused() {
				query := m.search.input.Value()
				if query != "" {
					return m, m.doSearch(query)
				}
			}
			return m, nil

		case "down", "j":
			if !m.search.input.Focused() && m.search.cursor < len(m.search.results)-1 {
				m.search.cursor++
			}

		case "up", "k":
			if !m.search.input.Focused() && m.search.cursor > 0 {
				m.search.cursor--
			}

		case "right":
			if !m.search.input.Focused() && len(m.search.results) > 0 {
				old := m.search.selected[m.search.cursor]
				m.search.selected[m.search.cursor] = m.search.qty
				_ = old
			}

		case "left":
			if !m.search.input.Focused() {
				delete(m.search.selected, m.search.cursor)
			}

		case "+":
			if !m.search.input.Focused() {
				m.search.qty++
			}

		case "-":
			if !m.search.input.Focused() && m.search.qty > 1 {
				m.search.qty--
			}

		case "p":
			if !m.search.input.Focused() {
				if len(m.search.selected) == 0 {
					m.search.err = "No articles selected. Use → to select."
					return m, nil
				}
				if len(m.templates) == 0 {
					m.search.err = "No templates found. Add .zpl files to ~/.label-tui/templates/"
					return m, nil
				}
				m.current = screenSelectTemplate
				m.tmplSelect = newTmplSelectModel()
				return m, nil
			}

		case "tab":
			if m.search.input.Focused() {
				m.search.input.Blur()
			} else {
				m.search.input.Focus()
			}
		}
	}

	var cmd tea.Cmd
	m.search.input, cmd = m.search.input.Update(msg)
	return m, cmd
}

type searchResultsMsg struct {
	results []core.Article
	err     string
}

func (m *Model) doSearch(query string) tea.Cmd {
	return func() tea.Msg {
		q := strings.ReplaceAll(query, "'", "''")
		qUpper := strings.ToUpper(q)
		filter := fmt.Sprintf("contains(toupper(ItemCode),'%s') or contains(toupper(ItemName),'%s') or contains(toupper(BarCode),'%s')", qUpper, qUpper, qUpper)
		articles, err := m.sapClient.QueryArticles(filter)
		if err != nil {
			return searchResultsMsg{err: err.Error()}
		}
		return searchResultsMsg{results: articles}
	}
}

func (m *Model) viewSearch() string {
	var body string

	body += SectionStyle.Render("Search Articles") + "\n\n"

	searchInput := InputStyle
	if m.search.input.Focused() {
		searchInput = InputFocusedStyle
	}
	body += searchInput.Render(m.search.input.View()) + "\n"

	if m.search.loading {
		body += "\n" + SpinnerStyle.Render("Searching SAP B1...") + "\n"
		return body
	}

	if m.search.err != "" {
		body += "\n" + ErrorStyle.Render(m.search.err) + "\n"
	}

	if len(m.search.results) > 0 {
		selectedCount := len(m.search.selected)
		statusLine := fmt.Sprintf("Qty: %d  •  Selected: %d", m.search.qty, selectedCount)
		body += "\n" + InfoStyle.Render(statusLine) + "\n"

		for i, a := range m.search.results {
			_, checked := m.search.selected[i]

			line := fmt.Sprintf("%-20s %-30s $%8.2f", a.ItemCode, truncate(a.Description, 28), a.Price)
			if checked {
				line += fmt.Sprintf("  x%d ✓", m.search.selected[i])
				body += ArticleCheckedStyle.Render(line) + "\n"
			} else if i == m.search.cursor && !m.search.input.Focused() {
				body += ArticleSelectedStyle.Render("▸ "+line) + "\n"
			} else {
				body += ArticleItemStyle.Render("  "+line) + "\n"
			}
		}
	} else {
		body += "\n" + DimmedStyle.Render("No results yet. Enter a search term and press enter.") + "\n"
	}

	body += "\n" + RenderHelp(
		"tab focus list",
		"↑/↓ navigate",
		"→ select",
		"← deselect",
		"+/- qty",
		"p print",
	)

	return body
}

type tmplSelectModel struct {
	cursor int
	err    string
}

func newTmplSelectModel() tmplSelectModel {
	return tmplSelectModel{}
}

func (m *Model) updateTemplateSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down", "j":
			if m.tmplSelect.cursor < len(m.templates)-1 {
				m.tmplSelect.cursor++
			}
		case "up", "k":
			if m.tmplSelect.cursor > 0 {
				m.tmplSelect.cursor--
			}
		case "enter":
			if len(m.templates) > 0 {
				tmpl := m.templates[m.tmplSelect.cursor]
				m.current = screenPreview
				m.preview = newPreviewModel(m.search.selected, m.search.results, tmpl)
			}
		}
	}
	return m, nil
}

func (m *Model) viewTemplateSelect() string {
	body := SectionStyle.Render("Select Label Template") + "\n\n"

	for i, t := range m.templates {
		line := "  " + t
		if i == m.tmplSelect.cursor {
			body += MenuSelectedStyle.Render("▸ "+t) + "\n"
		} else {
			body += MenuItemStyle.Render("  "+t) + "\n"
		}
		_ = line
	}

	body += "\n" + RenderHelp("↑/↓ navigate", "enter select", "esc back")
	return body
}

type previewModel struct {
	selected  map[int]int
	articles  []core.Article
	tmplName  string
	rawZPL    string
	err       string
	printing  bool
	printed   bool
}

func newPreviewModel(selected map[int]int, articles []core.Article, tmplName string) previewModel {
	m := previewModel{
		selected: selected,
		articles: articles,
		tmplName: tmplName,
	}
	raw, err := core.LoadTemplate(tmplName)
	if err != nil {
		m.err = err.Error()
	} else {
		m.rawZPL = raw
	}
	return m
}

func (m *Model) updatePreview(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.preview.printing {
		return m, nil
	}

	switch msg := msg.(type) {
	case printResultMsg:
		m.preview.printing = false
		if msg.err != "" {
			m.preview.err = msg.err
		} else {
			m.preview.printed = true
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, m.doPrint()
		case "r":
			m.preview.printed = false
			m.preview.err = ""
			return m, nil
		}
	}

	return m, nil
}

type printResultMsg struct {
	err string
}

func (m *Model) doPrint() tea.Cmd {
	return func() tea.Msg {
		zpl, err := core.LoadTemplate(m.preview.tmplName)
		if err != nil {
			return printResultMsg{err: err.Error()}
		}

		if m.settings.USBPort == "" {
			return printResultMsg{err: "USB port not configured. Go to Settings first."}
		}

		var fullZPL string
		for idx, qty := range m.preview.selected {
			article := m.preview.articles[idx]
			fullZPL += core.RenderTemplate(zpl, article, qty)
		}

		if err := printers.PrintZPL(m.settings.USBPort, fullZPL); err != nil {
			return printResultMsg{err: err.Error()}
		}

		return printResultMsg{}
	}
}

func (m *Model) viewPreview() string {
	body := SectionStyle.Render("Preview & Print") + "\n\n"

	info := fmt.Sprintf("Template: %s  •  Articles: %d  •  Total labels: %d",
		m.preview.tmplName, len(m.preview.selected), totalLabels(m.preview.selected))
	body += InfoStyle.Render(info) + "\n\n"

	body += PreviewBoxStyle.Render(func() string {
		var s string
		s += SectionStyle.Render("Selected Articles") + "\n"
		for idx, qty := range m.preview.selected {
			a := m.preview.articles[idx]
			s += fmt.Sprintf("  %s  %s  x%d  $%.2f\n", a.ItemCode, truncate(a.Description, 30), qty, a.Price)
		}
		return s
	}()) + "\n"

	if m.preview.rawZPL != "" {
		body += PreviewBoxStyle.Render(func() string {
			zpl := m.preview.rawZPL
			if len(zpl) > 400 {
				zpl = zpl[:400]
			}
			return SectionStyle.Render("ZPL Output") + "\n" + zpl
		}()) + "\n"
	}

	if m.preview.err != "" {
		body += ErrorStyle.Render("Error: "+m.preview.err) + "\n"
	} else if m.preview.printing {
		body += "\n" + SpinnerStyle.Render("Sending to printer...") + "\n"
	} else if m.preview.printed {
		body += "\n" + SuccessStyle.Render("✓ Labels sent to printer successfully!") + "\n"
		body += "\n" + RenderHelp("r print again", "esc back")
	} else {
		body += "\n" + RenderHelp("enter print", "esc back")
	}

	return body
}

type settingsModel struct {
	inputs       []textinput.Model
	focused      int
	saved        bool
	err          string
	ports        []string
	width, height int
}

func newSettingsModel(s *core.Settings) settingsModel {
	inputs := make([]textinput.Model, 4)

	inputs[0] = textinput.New()
	inputs[0].Placeholder = "SBODemoCL"
	inputs[0].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorPrimary))
	if s.CompanyDB != "" {
		inputs[0].SetValue(s.CompanyDB)
	}
	inputs[0].Focus()

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "https://your-server:50000/b1s/v1"
	inputs[1].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorPrimary))
	if s.SAPServiceLayerURL != "" {
		inputs[1].SetValue(s.SAPServiceLayerURL)
	}

	inputs[2] = textinput.New()
	inputs[2].Placeholder = "/dev/usb/lp0  or  COM3"
	inputs[2].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorPrimary))
	if s.USBPort != "" {
		inputs[2].SetValue(s.USBPort)
	}

	inputs[3] = textinput.New()
	inputs[3].Placeholder = "1"
	inputs[3].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorPrimary))
	if s.PriceList > 0 {
		inputs[3].SetValue(fmt.Sprintf("%d", s.PriceList))
	}

	m := settingsModel{inputs: inputs}
	m.detectPorts()
	return m
}

func (sm *settingsModel) detectPorts() {
	ports, err := printers.ListUSBPorts()
	if err == nil {
		sm.ports = ports
	}
}

func (m *Model) updateSettings(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			m.settingsScr.focused = (m.settingsScr.focused + 1) % len(m.settingsScr.inputs)
			return m, m.focusSettingsInput()
		case "shift+tab", "up":
			m.settingsScr.focused = (m.settingsScr.focused - 1 + len(m.settingsScr.inputs)) % len(m.settingsScr.inputs)
			return m, m.focusSettingsInput()
		case "enter":
			return m, m.saveSettings()
		}
	}

	cmd := m.updateSettingsInputs(msg)
	return m, cmd
}

func (m *Model) focusSettingsInput() tea.Cmd {
	var cmds []tea.Cmd
	for i := range m.settingsScr.inputs {
		if i == m.settingsScr.focused {
			cmds = append(cmds, m.settingsScr.inputs[i].Focus())
		} else {
			m.settingsScr.inputs[i].Blur()
		}
	}
	return tea.Batch(cmds...)
}

func (m *Model) updateSettingsInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.settingsScr.inputs))
	for i := range m.settingsScr.inputs {
		m.settingsScr.inputs[i], cmds[i] = m.settingsScr.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m *Model) saveSettings() tea.Cmd {
	return func() tea.Msg {
		m.settings.CompanyDB = m.settingsScr.inputs[0].Value()
		m.settings.SAPServiceLayerURL = m.settingsScr.inputs[1].Value()
		m.settings.USBPort = m.settingsScr.inputs[2].Value()
		fmt.Sscanf(m.settingsScr.inputs[3].Value(), "%d", &m.settings.PriceList)

		if m.settings.SAPServiceLayerURL != "" {
			m.sapClient = core.NewSapClient(m.settings.SAPServiceLayerURL, m.settings.PriceList)
		}

		if err := core.SaveSettings(m.settings); err != nil {
			return settingsSavedMsg{err: err.Error()}
		}
		return settingsSavedMsg{saved: true}
	}
}

type settingsSavedMsg struct {
	saved bool
	err   string
}

func (m *Model) viewSettings() string {
	body := SectionStyle.Render("Configuration") + "\n\n"

	labels := []string{"CompanyDB", "SAP Service Layer URL", "USB Printer Port", "Price List (number)"}
	for i, input := range m.settingsScr.inputs {
		body += DimmedStyle.Render(labels[i]) + "\n"
		style := InputStyle
		if i == m.settingsScr.focused {
			style = InputFocusedStyle
		}
		body += style.Render(input.View()) + "\n\n"
	}

	if len(m.settingsScr.ports) > 0 {
		body += SectionStyle.Render("Detected USB Ports") + "\n"
		for _, p := range m.settingsScr.ports {
			body += DimmedStyle.Render("  • "+p) + "\n"
		}
		body += "\n"
	} else {
		body += DimmedStyle.Render("No USB ports detected. Connect your printer.") + "\n\n"
	}

	if m.settingsScr.err != "" {
		body += ErrorStyle.Render("Error: "+m.settingsScr.err) + "\n"
	} else if m.settingsScr.saved {
		body += SuccessStyle.Render("✓ Settings saved successfully!") + "\n"
	}

	body += RenderHelp("tab/shift+tab navigate", "enter save", "esc back")
	return body
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

func totalLabels(selected map[int]int) int {
	total := 0
	for _, qty := range selected {
		total += qty
	}
	return total
}
