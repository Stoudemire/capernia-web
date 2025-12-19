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
        }

        AdminNewsTmplData struct {
                Common CommonTmplData
                NewsList []TNews
                EditingNews *TNews
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
                "title": strings.Title,
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
        
        ExecuteTemplate(Context.Writer, "highscores.tmpl",
                HighscoresTmplData{
                        Common: GetCommonTmplData("Highscores", Context.AccountID),
                        Highscores:        GetHighscores(Skill, Vocation),
                        CurrentSkill:      Skill,
                        CurrentSkillDisplay: skillDisp,
                        CurrentVocation:   Vocation,
                })
}

func RenderNews(Context *THttpRequestContext) {
        news, err := GetAllNews()
        if err != nil {
                g_LogErr.Printf("Failed to get news: %v", err)
                news = []TNews{}
        }

        ExecuteTemplate(Context.Writer, "news.tmpl",
                NewsTmplData{
                        Common: GetCommonTmplData("News", Context.AccountID),
                        NewsList: news,
                })
}

func RenderAdminNews(Context *THttpRequestContext) {
        news, err := GetAllNews()
        if err != nil {
                g_LogErr.Printf("Failed to get news: %v", err)
                news = []TNews{}
        }

        ExecuteTemplate(Context.Writer, "admin_news.tmpl",
                AdminNewsTmplData{
                        Common: GetCommonTmplData("Admin News", Context.AccountID),
                        NewsList: news,
                        EditingNews: nil,
                })
}

func RenderAdminNewsEdit(Context *THttpRequestContext, newsID int) {
        news, err := GetAllNews()
        if err != nil {
                g_LogErr.Printf("Failed to get news: %v", err)
                news = []TNews{}
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
                })
}
