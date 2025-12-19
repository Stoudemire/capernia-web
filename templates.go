package main

import (
        "fmt"
        "html/template"
        "io"
        "net/http"
        "strconv"
        "strings"
)

type (
        CommonTmplData struct {
                Title         string
                AccountID     int
                TotalPlayers  int
                ServerOnline  bool
                LastStartup   int
        }

        GenericTmplData struct {
                Common CommonTmplData
        }

        NewsTmplData struct {
                Common CommonTmplData
                NewsList []TNews
                CurrentPage int
                TotalPages int
                TotalNews int
        }

        AdminNewsTmplData struct {
                Common CommonTmplData
                NewsList []TNews
                EditingNews *TNews
                CurrentPage int
                TotalPages int
                TotalNews int
        }

        AccountTmplData struct {
                Common  CommonTmplData
                Account *TAccountSummary
        }

        CharacterTmplData struct {
                Common    CommonTmplData
                Character *TCharacterProfile
        }

        KillStatisticsTmplData struct {
                Common         CommonTmplData
                World          *TWorld
                KillStatistics []TKillStatistics
        }

        WorldTmplData struct {
                Common           CommonTmplData
                World            *TWorld
                OnlineCharacters []TOnlineCharacter
        }

        WorldListTmplData struct {
                Common CommonTmplData
                Worlds []TWorld
        }

        MessageTmplData struct {
                Common  CommonTmplData
                Heading string
                Message string
        }

        THighscore struct {
                CharacterName string
                Profession    string
                Level         int
                SkillName     string
                SkillValue    int
        }

        HighscoresTmplData struct {
                Common            CommonTmplData
                Highscores        []THighscore
                CurrentSkill      string
                CurrentSkillName  string
                CurrentSkillDisplay string
                CurrentVocation   string
                CurrentPage       int
                TotalPages        int
                TotalHighscores   int
        }

        NewsArchiveTmplData struct {
                Common      CommonTmplData
                SearchNews  []TNews
                HasResults  bool
                CurrentPage int
                TotalPages  int
                TotalNews   int
                FromDay     string
                FromMonth   string
                FromYear    string
                ToDay       string
                ToMonth     string
                ToYear      string
        }

        THouse struct {
                HouseID     int
                Name        string
                Rent        int
                Description string
                Size        int
                Town        string
                GuildHouse  bool
                Status      string
                Owner       string
                PaidUntil   int
        }

        HousesTmplData struct {
                Common      CommonTmplData
                Houses      []THouse
                SelectedTown string
                SelectedType int
                SelectedStatus int
                Towns       []string
        }

        HouseDetailTmplData struct {
                Common CommonTmplData
                House  *THouse
        }

        TGuild struct {
                GuildID   int
                Name      string
                Leader    string
                Created   int
                MemberCount int
        }

        TGuildMember struct {
                CharacterName string
                Rank          string
                Title         string
                Joined        int
                Level         int
                Profession    string
        }

        GuildsTmplData struct {
                Common CommonTmplData
                Guilds []TGuild
        }

        GuildDetailTmplData struct {
                Common  CommonTmplData
                Guild   *TGuild
                Members []TGuildMember
        }
)

var (
        g_Templates *template.Template
)

func InitTemplates() bool {
        var Err error

        CustomFuncs := template.FuncMap{
                "FormatTimestamp": FormatTimestamp,
                "FormatDurationSince": FormatDurationSince,
                "add": func(a, b int) int { return a + b },
                "sub": func(a, b int) int { return a - b },
                "mul": func(a, b int) int { return a * b },
                "until": func(n int) []int {
                        result := make([]int, n)
                        for i := 0; i < n; i++ {
                                result[i] = i
                        }
                        return result
                },
                "title": strings.Title,
                "HTML": func(s string) template.HTML { return template.HTML(s) },
        }

        g_Templates, Err = template.New("").Funcs(CustomFuncs).ParseGlob("templates/*.tmpl")
        if Err != nil {
                g_LogErr.Printf("Failed to parse templates: %v", Err)
                return false
        }
        return true
}

func ExitTemplates() {
        g_Templates = nil
}

func ExecuteTemplate(Writer io.Writer, FileName string, Data any) {
        Err := g_Templates.ExecuteTemplate(Writer, FileName, Data)
        if Err != nil {
                g_LogErr.Printf("Failed to execute template \"%v\": %v", FileName, Err)
        }
}

func GetCommonTmplData(Title string, AccountID int) CommonTmplData {
        Worlds := GetWorlds()
        TotalPlayers := 0
        ServerOnline := false
        LastStartup := 0
        
        for _, World := range Worlds {
                TotalPlayers += World.NumPlayers
                if World.LastStartup > World.LastShutdown {
                        ServerOnline = true
                }
                if World.LastStartup > LastStartup {
                        LastStartup = World.LastStartup
                }
        }
        
        return CommonTmplData{
                Title:         Title,
                AccountID:     AccountID,
                TotalPlayers:  TotalPlayers,
                ServerOnline:  ServerOnline,
                LastStartup:   LastStartup,
        }
}

func RenderRequestError(Context *THttpRequestContext, Status int) {
        StatusText := http.StatusText(Status)
        ExecuteTemplate(Context.Writer, "message.tmpl",
                MessageTmplData{
                        Common:  GetCommonTmplData(StatusText, Context.AccountID),
                        Heading: strconv.Itoa(Status),
                        Message: StatusText,
                })
}

func RenderMessage(Context *THttpRequestContext, Heading string, Message string) {
        ExecuteTemplate(Context.Writer, "message.tmpl",
                MessageTmplData{
                        Common:  GetCommonTmplData(Heading, Context.AccountID),
                        Heading: Heading,
                        Message: Message,
                })
}

func RenderAccountSummary(Context *THttpRequestContext) {
        Data := AccountTmplData{
                Common: GetCommonTmplData("Account Summary", Context.AccountID),
                Account: nil,
        }

        Result, Account := GetAccountSummary(Context.AccountID)
        if Result == 0 {
                Data.Account = &Account
        }

        ExecuteTemplate(Context.Writer, "account_summary.tmpl", Data)
}

func RenderAccountLogin(Context *THttpRequestContext) {
        ExecuteTemplate(Context.Writer, "account_login.tmpl",
                GenericTmplData{
                        Common: GetCommonTmplData("Login", Context.AccountID),
                })
}

func RenderAccountCreate(Context *THttpRequestContext) {
        ExecuteTemplate(Context.Writer, "account_create.tmpl",
                GenericTmplData{
                        Common: GetCommonTmplData("Create Account", Context.AccountID),
                })
}

func RenderAccountRecover(Context *THttpRequestContext) {
        ExecuteTemplate(Context.Writer, "account_recover.tmpl",
                GenericTmplData{
                        Common: GetCommonTmplData("Recover Account", Context.AccountID),
                })
}

func RenderCharacterCreate(Context *THttpRequestContext) {
        ExecuteTemplate(Context.Writer, "character_create.tmpl",
                WorldListTmplData{
                        Common: GetCommonTmplData("Create Character", Context.AccountID),
                        Worlds: GetWorlds(),
                })
}

func RenderCharacterProfile(Context *THttpRequestContext, Character *TCharacterProfile) {
        Title := "Search Character"
        if Character != nil {
                Title = fmt.Sprintf("%v's Profile", Character.Name)
        }

        ExecuteTemplate(Context.Writer, "character_profile.tmpl",
                CharacterTmplData{
                        Common: GetCommonTmplData(Title, Context.AccountID),
                        Character: Character,
                })
}

func RenderKillStatisticsList(Context *THttpRequestContext) {
        ExecuteTemplate(Context.Writer, "killstatistics_list.tmpl",
                WorldListTmplData{
                        Common: GetCommonTmplData("Kill Statistics", Context.AccountID),
                        Worlds: GetWorlds(),
                })
}

func RenderKillStatistics(Context *THttpRequestContext, WorldName string) {
        ExecuteTemplate(Context.Writer, "killstatistics.tmpl",
                KillStatisticsTmplData{
                        Common: GetCommonTmplData(fmt.Sprintf("Kill Statistics - %v", WorldName), Context.AccountID),
                        World:          GetWorld(WorldName),
                        KillStatistics: GetKillStatistics(WorldName),
                })
}

func RenderWorldList(Context *THttpRequestContext) {
        ExecuteTemplate(Context.Writer, "world_list.tmpl",
                WorldListTmplData{
                        Common: GetCommonTmplData("Worlds", Context.AccountID),
                        Worlds: GetWorlds(),
                })
}

func RenderWorldInfo(Context *THttpRequestContext, WorldName string) {
        ExecuteTemplate(Context.Writer, "world_info.tmpl",
                WorldTmplData{
                        Common: GetCommonTmplData("Worlds", Context.AccountID),
                        World:            GetWorld(WorldName),
                        OnlineCharacters: GetOnlineCharacters(WorldName),
                })
}

func RenderHighscores(Context *THttpRequestContext, Skill string, Vocation string) {
        skillDisplay := map[string]string{
                "level": "Level",
                "magic": "Magic",
                "fist": "Fist Fighting",
                "club": "Club Fighting",
                "sword": "Sword Fighting",
                "axe": "Axe Fighting",
                "distance": "Distance Fighting",
                "shielding": "Shielding",
                "fishing": "Fishing",
        }
        
        skillDisp := "Level"
        if disp, ok := skillDisplay[Skill]; ok {
                skillDisp = disp
        }
        
        pageStr := Context.Request.URL.Query().Get("page")
        page := 1
        if pageStr != "" {
                p, err := strconv.Atoi(pageStr)
                if err == nil && p > 0 {
                        page = p
                }
        }

        itemsPerPage := 50
        allHighscores := GetHighscores(Skill, Vocation)
        totalHighscores := len(allHighscores)
        totalPages := (totalHighscores + itemsPerPage - 1) / itemsPerPage
        if totalPages < 1 {
                totalPages = 1
        }

        startIdx := (page - 1) * itemsPerPage
        endIdx := startIdx + itemsPerPage
        if startIdx >= totalHighscores {
                startIdx = 0
                page = 1
        }
        if endIdx > totalHighscores {
                endIdx = totalHighscores
        }

        paginatedHighscores := allHighscores[startIdx:endIdx]
        
        ExecuteTemplate(Context.Writer, "highscores.tmpl",
                HighscoresTmplData{
                        Common: GetCommonTmplData("Highscores", Context.AccountID),
                        Highscores:        paginatedHighscores,
                        CurrentSkill:      Skill,
                        CurrentSkillDisplay: skillDisp,
                        CurrentVocation:   Vocation,
                        CurrentPage:       page,
                        TotalPages:        totalPages,
                        TotalHighscores:   totalHighscores,
                })
}

func RenderNews(Context *THttpRequestContext) {
        pageStr := Context.Request.URL.Query().Get("page")
        page := 1
        if pageStr != "" {
                p, err := strconv.Atoi(pageStr)
                if err == nil && p > 0 {
                        page = p
                }
        }

        itemsPerPage := 3
        news, err := GetNewsPaginated(page, itemsPerPage)
        if err != nil {
                g_LogErr.Printf("Failed to get news: %v", err)
                news = []TNews{}
        }

        totalNews, err := GetTotalNewsCount()
        if err != nil {
                g_LogErr.Printf("Failed to get total news count: %v", err)
                totalNews = 0
        }

        totalPages := (totalNews + itemsPerPage - 1) / itemsPerPage
        if totalPages < 1 {
                totalPages = 1
        }

        ExecuteTemplate(Context.Writer, "news.tmpl",
                NewsTmplData{
                        Common: GetCommonTmplData("News", Context.AccountID),
                        NewsList: news,
                        CurrentPage: page,
                        TotalPages: totalPages,
                        TotalNews: totalNews,
                })
}

func RenderNewsArchive(Context *THttpRequestContext, Data *NewsArchiveTmplData) {
        ExecuteTemplate(Context.Writer, "news_archive.tmpl", Data)
}

func RenderAdminNews(Context *THttpRequestContext) {
        pageStr := Context.Request.URL.Query().Get("page")
        page := 1
        if pageStr != "" {
                p, err := strconv.Atoi(pageStr)
                if err == nil && p > 0 {
                        page = p
                }
        }

        itemsPerPage := 3
        news, err := GetNewsPaginated(page, itemsPerPage)
        if err != nil {
                g_LogErr.Printf("Failed to get news: %v", err)
                news = []TNews{}
        }

        totalNews, err := GetTotalNewsCount()
        if err != nil {
                g_LogErr.Printf("Failed to get total news count: %v", err)
                totalNews = 0
        }

        totalPages := (totalNews + itemsPerPage - 1) / itemsPerPage
        if totalPages < 1 {
                totalPages = 1
        }

        ExecuteTemplate(Context.Writer, "admin_news.tmpl",
                AdminNewsTmplData{
                        Common: GetCommonTmplData("Admin News", Context.AccountID),
                        NewsList: news,
                        EditingNews: nil,
                        CurrentPage: page,
                        TotalPages: totalPages,
                        TotalNews: totalNews,
                })
}

func RenderAdminNewsEdit(Context *THttpRequestContext, newsID int) {
        pageStr := Context.Request.URL.Query().Get("page")
        page := 1
        if pageStr != "" {
                p, err := strconv.Atoi(pageStr)
                if err == nil && p > 0 {
                        page = p
                }
        }

        itemsPerPage := 3
        news, err := GetNewsPaginated(page, itemsPerPage)
        if err != nil {
                g_LogErr.Printf("Failed to get news: %v", err)
                news = []TNews{}
        }

        totalNews, err := GetTotalNewsCount()
        if err != nil {
                g_LogErr.Printf("Failed to get total news count: %v", err)
                totalNews = 0
        }

        totalPages := (totalNews + itemsPerPage - 1) / itemsPerPage
        if totalPages < 1 {
                totalPages = 1
        }

        editingNews, err := GetNewsById(newsID)
        if err != nil {
                g_LogErr.Printf("Failed to get news by id: %v", err)
        }

        ExecuteTemplate(Context.Writer, "admin_news.tmpl",
                AdminNewsTmplData{
                        Common: GetCommonTmplData("Admin News", Context.AccountID),
                        NewsList: news,
                        EditingNews: editingNews,
                        CurrentPage: page,
                        TotalPages: totalPages,
                        TotalNews: totalNews,
                })
}

func RenderDownloadClient(Context *THttpRequestContext) {
        ExecuteTemplate(Context.Writer, "download_client.tmpl",
                GenericTmplData{
                        Common: GetCommonTmplData("Download Client", Context.AccountID),
                })
}

func RenderHouses(Context *THttpRequestContext, Houses []THouse, SelectedTown string, SelectedType int, SelectedStatus int, Towns []string) {
        ExecuteTemplate(Context.Writer, "houses.tmpl",
                HousesTmplData{
                        Common: GetCommonTmplData("Houses", Context.AccountID),
                        Houses: Houses,
                        SelectedTown: SelectedTown,
                        SelectedType: SelectedType,
                        SelectedStatus: SelectedStatus,
                        Towns: Towns,
                })
}

func RenderHouseDetail(Context *THttpRequestContext, House *THouse) {
        ExecuteTemplate(Context.Writer, "house_detail.tmpl",
                HouseDetailTmplData{
                        Common: GetCommonTmplData("House", Context.AccountID),
                        House: House,
                })
}

func RenderGuilds(Context *THttpRequestContext, Guilds []TGuild) {
        ExecuteTemplate(Context.Writer, "guilds.tmpl",
                GuildsTmplData{
                        Common: GetCommonTmplData("Guilds", Context.AccountID),
                        Guilds: Guilds,
                })
}

func RenderGuildDetail(Context *THttpRequestContext, Guild *TGuild, Members []TGuildMember) {
        ExecuteTemplate(Context.Writer, "guild_detail.tmpl",
                GuildDetailTmplData{
                        Common: GetCommonTmplData("Guild", Context.AccountID),
                        Guild: Guild,
                        Members: Members,
                })
}
