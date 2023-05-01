package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/RunnersRevival/outrun/bgtasks"
	"github.com/RunnersRevival/outrun/config"
	"github.com/RunnersRevival/outrun/config/campaignconf"
	"github.com/RunnersRevival/outrun/config/eventconf"
	"github.com/RunnersRevival/outrun/config/gameconf"
	"github.com/RunnersRevival/outrun/config/infoconf"
	"github.com/RunnersRevival/outrun/cryption"
	"github.com/RunnersRevival/outrun/inforeporters"
	"github.com/RunnersRevival/outrun/meta"
	"github.com/RunnersRevival/outrun/muxhandlers"
	"github.com/RunnersRevival/outrun/muxhandlers/muxobj"
	"github.com/RunnersRevival/outrun/orpc"
	"github.com/gorilla/mux"
)

const UNKNOWN_REQUEST_DIRECTORY = "logging/unknown_requests/"

var (
	LogExecutionTime = true
)

func OutputUnknownRequest(w http.ResponseWriter, r *http.Request) {
	recv, err := cryption.GetReceivedMessage(r)

        if err != nil {
                panic(err)
        }
	// make a new logging path
	timeStr := strconv.Itoa(int(time.Now().Unix()))
	os.MkdirAll(UNKNOWN_REQUEST_DIRECTORY, 0644)
	normalizedReq := strings.ReplaceAll(r.URL.Path, "/", "-")
	path := UNKNOWN_REQUEST_DIRECTORY + normalizedReq + "_" + timeStr + ".txt"
	err := ioutil.WriteFile(path, recv, 0644)
	if err != nil {
		log.Println("[OUT] UNABLE TO WRITE UNKNOWN REQUEST: " + err.Error())
		w.Write([]byte(""))
		return
	}
	log.Println("[OUT] !!!!!!!!!!!! Unknown request, output to " + path)
	w.Write([]byte(""))
}

func removePrependingSlashes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for len(r.URL.Path) != 0 && string(r.URL.Path[0]) == "/" {
			r.URL.Path = r.URL.Path[1:]
		}
		r.URL.Path = "/" + r.URL.Path
		next.ServeHTTP(w, r)
	})
}

func Generate204(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// return Not Found for the favicon; no favicon is intended
func FaviconResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Not Found"))
}

// Return "OK" for checking if the Outrun instance is alive (intended for uptime monitors)
func GenericRootResponse(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func checkArgs() bool {
	// TODO: _VERY_ dirty command line argument checking. This should be
	// changed into something more robust and less hacky!
	args := os.Args[1:] // drop executable
	amt := len(args)
	if amt >= 1 {
		if args[0] == "--version" {
			fmt.Printf("Outrun for Revival %s (self-hostable edition)\n", meta.Version)
			return true
		}
		fmt.Println("Unknown given arguments")
		return true
	}
	return false
}

func main() {
	end := checkArgs()
	if end {
		return
	}
	rand.Seed(time.Now().UTC().UnixNano())
	log.Println("[INFO] Outrun for Revival (self-hostable edition) is starting up...")

	err := config.Parse("config.json")
	if err != nil {
		log.Printf("[INFO] Failure loading config file config.json (%s), using defaults\n", err)
	} else {
		log.Println("[INFO] Config file (config.json) loaded")
	}

	err = eventconf.Parse(config.CFile.EventConfigFilename)
	if err != nil {
		if !config.CFile.SilenceEventConfigErrors {
			log.Printf("[INFO] Failure loading event config file %s (%s), using defaults\n", config.CFile.EventConfigFilename, err)
		}
	} else {
		log.Printf("[INFO] Event config file (%s) loaded\n", config.CFile.EventConfigFilename)
	}

	err = infoconf.Parse(config.CFile.InfoConfigFilename)
	if err != nil {
		if !config.CFile.SilenceInfoConfigErrors {
			log.Printf("[INFO] Failure loading info config file %s (%s), using defaults\n", config.CFile.InfoConfigFilename, err)
		}
	} else {
		log.Printf("[INFO] Info config file (%s) loaded\n", config.CFile.InfoConfigFilename)
	}

	err = gameconf.Parse(config.CFile.GameConfigFilename)
	if err != nil {
		if !config.CFile.SilenceGameConfigErrors {
			log.Printf("[INFO] Failure loading game config file %s (%s), using defaults\n", config.CFile.GameConfigFilename, err)
		}
	} else {
		log.Printf("[INFO] Game config file (%s) loaded\n", config.CFile.GameConfigFilename)
	}

	err = campaignconf.Parse(config.CFile.CampaignConfigFilename)
	if err != nil {
		if !config.CFile.SilenceCampaignConfigErrors {
			log.Printf("[INFO] Failure loading campaign config file %s (%s), using defaults\n", config.CFile.CampaignConfigFilename, err)
		}
	} else {
		log.Printf("[INFO] Campaign config file (%s) loaded\n", config.CFile.CampaignConfigFilename)
	}

	if config.CFile.EnableRPC {
		orpc.Start()
	}

	h := muxobj.Handle
	router := mux.NewRouter()
	router.StrictSlash(true)
	LogExecutionTime = config.CFile.DoTimeLogging
	prefix := config.CFile.EndpointPrefix

	router.HandleFunc("/", GenericRootResponse)
	router.HandleFunc("/favicon.ico", FaviconResponse)
	
	// Login
	router.HandleFunc(prefix+"/Login/login/", h(muxhandlers.Login, LogExecutionTime))
	router.HandleFunc(prefix+"/Sgn/sendApollo/", h(muxhandlers.SendApollo, LogExecutionTime))
	router.HandleFunc(prefix+"/Login/getVariousParameter/", h(muxhandlers.GetVariousParameter, LogExecutionTime))
	router.HandleFunc(prefix+"/Player/getPlayerState/", h(muxhandlers.GetPlayerState, LogExecutionTime))
	router.HandleFunc(prefix+"/Player/getCharacterState/", h(muxhandlers.GetCharacterState, LogExecutionTime))
	router.HandleFunc(prefix+"/Player/getChaoState/", h(muxhandlers.GetChaoState, LogExecutionTime))
	router.HandleFunc(prefix+"/Spin/getWheelOptions/", h(muxhandlers.GetWheelOptions, LogExecutionTime))
	router.HandleFunc(prefix+"/Game/getDailyChalData/", h(muxhandlers.GetDailyChallengeData, LogExecutionTime))
	router.HandleFunc(prefix+"/Message/getMessageList/", h(muxhandlers.GetMessageList, LogExecutionTime))
	router.HandleFunc(prefix+"/Store/getRedstarExchangeList/", h(muxhandlers.GetRedStarExchangeList, LogExecutionTime))
	router.HandleFunc(prefix+"/Game/getCostList/", h(muxhandlers.GetCostList, LogExecutionTime))
	router.HandleFunc(prefix+"/Event/getEventList/", h(muxhandlers.GetEventList, LogExecutionTime))
	router.HandleFunc(prefix+"/Game/getMileageData/", h(muxhandlers.GetMileageData, LogExecutionTime))
	router.HandleFunc(prefix+"/Game/getCampaignList/", h(muxhandlers.GetCampaignList, LogExecutionTime))
	router.HandleFunc(prefix+"/Chao/getChaoWheelOptions/", h(muxhandlers.GetChaoWheelOptions, LogExecutionTime))
	router.HandleFunc(prefix+"/Chao/getPrizeChaoWheelSpin/", h(muxhandlers.GetPrizeChaoWheelSpin, LogExecutionTime))
	router.HandleFunc(prefix+"/login/getInformation/", h(muxhandlers.GetInformation, LogExecutionTime))
	router.HandleFunc(prefix+"/Leaderboard/getWeeklyLeaderboardOptions/", h(muxhandlers.GetWeeklyLeaderboardOptions, LogExecutionTime))
	router.HandleFunc(prefix+"/Leaderboard/getLeagueData/", h(muxhandlers.GetLeagueData, LogExecutionTime))
	router.HandleFunc(prefix+"/Leaderboard/getWeeklyLeaderboardEntries/", h(muxhandlers.GetWeeklyLeaderboardEntries, LogExecutionTime))
	router.HandleFunc(prefix+"/Player/setUserName/", h(muxhandlers.SetUsername, LogExecutionTime))
	router.HandleFunc(prefix+"/login/getTicker/", h(muxhandlers.GetTicker, LogExecutionTime))
	router.HandleFunc(prefix+"/Login/loginBonus/", h(muxhandlers.LoginBonus, LogExecutionTime))
	router.HandleFunc(prefix+"/Login/getCountry/", h(muxhandlers.GetCountry, LogExecutionTime))
	router.HandleFunc(prefix+"/Option/userResult/", h(muxhandlers.GetOptionUserResult, LogExecutionTime))
	router.HandleFunc(prefix+"/Message/getMessage/", h(muxhandlers.GetMessage, LogExecutionTime))
	router.HandleFunc(prefix+"/Login/loginBonusSelect/", h(muxhandlers.LoginBonusSelect, LogExecutionTime))
	// Timed mode
	router.HandleFunc(prefix+"/Game/quickActStart/", h(muxhandlers.QuickActStart, LogExecutionTime))
	router.HandleFunc(prefix+"/Game/quickPostGameResults/", h(muxhandlers.QuickPostGameResults, LogExecutionTime))
	// Story mode
	router.HandleFunc(prefix+"/Game/actStart/", h(muxhandlers.ActStart, LogExecutionTime))
	router.HandleFunc(prefix+"/Game/getMileageReward/", h(muxhandlers.GetMileageReward, LogExecutionTime))
	// Retry
	router.HandleFunc(prefix+"/Game/actRetry/", h(muxhandlers.ActRetry, LogExecutionTime))
	// Gameplay
	router.HandleFunc(prefix+"/Game/getFreeItemList/", h(muxhandlers.GetFreeItemList, LogExecutionTime))
	router.HandleFunc(prefix+"/Game/postGameResults/", h(muxhandlers.PostGameResults, LogExecutionTime))
	// Misc.
	router.HandleFunc(prefix+"/Character/changeCharacter/", h(muxhandlers.ChangeCharacter, LogExecutionTime))
	router.HandleFunc(prefix+"/Character/upgradeCharacter/", h(muxhandlers.UpgradeCharacter, LogExecutionTime))
	router.HandleFunc(prefix+"/Chao/equipChao/", h(muxhandlers.EquipChao, LogExecutionTime))
	// Shop
	router.HandleFunc(prefix+"/Store/redstarExchange/", h(muxhandlers.RedStarExchange, LogExecutionTime))
	// Friends (Required for iOS?)
	router.HandleFunc(prefix+"/Friend/getFacebookIncentive/", h(muxhandlers.GetFacebookIncentive, LogExecutionTime))
	// Roulette
	router.HandleFunc(prefix+"/RaidbossSpin/getItemStockNum/", h(muxhandlers.GetItemStockNum, LogExecutionTime))
	router.HandleFunc(prefix+"/Spin/commitWheelSpin/", h(muxhandlers.CommitWheelSpin, LogExecutionTime))
	router.HandleFunc(prefix+"/Chao/commitChaoWheelSpin/", h(muxhandlers.CommitChaoWheelSpin, LogExecutionTime))
	// Character transactions
	router.HandleFunc(prefix+"/Character/unlockedCharacter/", h(muxhandlers.UnlockedCharacter, LogExecutionTime))
	// Migration
	router.HandleFunc(prefix+"/Login/getMigrationPassword/", h(muxhandlers.GetMigrationPassword, LogExecutionTime))
	router.HandleFunc(prefix+"/Login/migration/", h(muxhandlers.Migration, LogExecutionTime))

	// Server information
	if config.CFile.EnablePublicStats {
		router.HandleFunc("/outrunInfo/stats", inforeporters.Stats)
	}

	if config.CFile.LogUnknownRequests {
		router.PathPrefix("/").HandlerFunc(OutputUnknownRequest)
	}
	
	router.HandleFunc("/generate204", Generate204)

	go bgtasks.TouchAnalyticsDB()

	port := config.CFile.Port
	log.Printf("Starting server on port %s\n", port)
	panic(http.ListenAndServe(":"+port, removePrependingSlashes(router)))
}
