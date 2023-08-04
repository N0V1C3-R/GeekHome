package command_server

import (
	"WebHome/src/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type TimeConvertServer struct {
	BaseCommand
}

type TimestampConvertServer struct {
	BaseCommand
}

var timeZones = []string{
	"Africa/Abidjan",
	"Africa/Accra",
	"Africa/Addis_Ababa",
	"Africa/Algiers",
	"Africa/Asmara",
	"Africa/Asmera",
	"Africa/Bamako",
	"Africa/Bangui",
	"Africa/Banjul",
	"Africa/Bissau",
	"Africa/Blantyre",
	"Africa/Brazzaville",
	"Africa/Bujumbura",
	"Africa/Cairo",
	"Africa/Casablanca",
	"Africa/Ceuta",
	"Africa/Conakry",
	"Africa/Dakar",
	"Africa/Dar_es_Salaam",
	"Africa/Djibouti",
	"Africa/Douala",
	"Africa/El_Aaiun",
	"Africa/Freetown",
	"Africa/Gaborone",
	"Africa/Harare",
	"Africa/Johannesburg",
	"Africa/Juba",
	"Africa/Kampala",
	"Africa/Khartoum",
	"Africa/Kigali",
	"Africa/Kinshasa",
	"Africa/Lagos",
	"Africa/Libreville",
	"Africa/Lome",
	"Africa/Luanda",
	"Africa/Lubumbashi",
	"Africa/Lusaka",
	"Africa/Malabo",
	"Africa/Maputo",
	"Africa/Maseru",
	"Africa/Mbabane",
	"Africa/Mogadishu",
	"Africa/Monrovia",
	"Africa/Nairobi",
	"Africa/Ndjamena",
	"Africa/Niamey",
	"Africa/Nouakchott",
	"Africa/Ouagadougou",
	"Africa/Porto-Novo",
	"Africa/Sao_Tome",
	"Africa/Tripoli",
	"Africa/Tunis",
	"Africa/Windhoek",
	"America/Adak",
	"America/Anchorage",
	"America/Anguilla",
	"America/Antigua",
	"America/Araguaina",
	"America/Argentina/Buenos_Aires",
	"America/Argentina/Catamarca",
	"America/Argentina/Cordoba",
	"America/Argentina/Jujuy",
	"America/Argentina/La_Rioja",
	"America/Argentina/Mendoza",
	"America/Argentina/Rio_Gallegos",
	"America/Argentina/Salta",
	"America/Argentina/San_Juan",
	"America/Argentina/San_Luis",
	"America/Argentina/Tucuman",
	"America/Argentina/Ushuaia",
	"America/Aruba",
	"America/Asuncion",
	"America/Atikokan",
	"America/Bahia",
	"America/Bahia_Banderas",
	"America/Barbados",
	"America/Belem",
	"America/Belize",
	"America/Blanc-Sablon",
	"America/Boa_Vista",
	"America/Bogota",
	"America/Boise",
	"America/Cambridge_Bay",
	"America/Campo_Grande",
	"America/Cancun",
	"America/Caracas",
	"America/Cayenne",
	"America/Cayman",
	"America/Chicago",
	"America/Chihuahua",
	"America/Ciudad_Juarez",
	"America/Costa_Rica",
	"America/Creston",
	"America/Cuiaba",
	"America/Curacao",
	"America/Danmarkshavn",
	"America/Dawson",
	"America/Dawson_Creek",
	"America/Denver",
	"America/Detroit",
	"America/Dominica",
	"America/Edmonton",
	"America/Eirunepe",
	"America/El_Salvador",
	"America/Fort_Nelson",
	"America/Fortaleza",
	"America/Glace_Bay",
	"America/Godthab",
	"America/Goose_Bay",
	"America/Grand_Turk",
	"America/Grenada",
	"America/Guadeloupe",
	"America/Guatemala",
	"America/Guayaquil",
	"America/Guyana",
	"America/Halifax",
	"America/Havana",
	"America/Hermosillo",
	"America/Indiana/Indianapolis",
	"America/Indiana/Knox",
	"America/Indiana/Marengo",
	"America/Indiana/Petersburg",
	"America/Indiana/Tell_City",
	"America/Indiana/Vevay",
	"America/Indiana/Vincennes",
	"America/Indiana/Winamac",
	"America/Indianapolis",
	"America/Inuvik",
	"America/Iqaluit",
	"America/Jamaica",
	"America/Juneau",
	"America/Kentucky/Louisville",
	"America/Kentucky/Monticello",
	"America/Kralendijk",
	"America/La_Paz",
	"America/Lima",
	"America/Los_Angeles",
	"America/Lower_Princes",
	"America/Maceio",
	"America/Managua",
	"America/Manaus",
	"America/Marigot",
	"America/Martinique",
	"America/Matamoros",
	"America/Mazatlan",
	"America/Menominee",
	"America/Merida",
	"America/Metlakatla",
	"America/Mexico_City",
	"America/Miquelon",
	"America/Moncton",
	"America/Monterrey",
	"America/Montevideo",
	"America/Montreal",
	"America/Montserrat",
	"America/Nassau",
	"America/New_York",
	"America/Nome",
	"America/Noronha",
	"America/North_Dakota/Beulah",
	"America/North_Dakota/Center",
	"America/North_Dakota/New_Salem",
	"America/Nuuk",
	"America/Ojinaga",
	"America/Panama",
	"America/Paramaribo",
	"America/Phoenix",
	"America/Port-au-Prince",
	"America/Port_of_Spain",
	"America/Porto_Velho",
	"America/Puerto_Rico",
	"America/Punta_Arenas",
	"America/Rankin_Inlet",
	"America/Recife",
	"America/Regina",
	"America/Resolute",
	"America/Rio_Branco",
	"America/Santarem",
	"America/Santiago",
	"America/Santo_Domingo",
	"America/Sao_Paulo",
	"America/Scoresbysund",
	"America/Sitka",
	"America/St_Barthelemy",
	"America/St_Johns",
	"America/St_Kitts",
	"America/St_Lucia",
	"America/St_Thomas",
	"America/St_Vincent",
	"America/Swift_Current",
	"America/Tegucigalpa",
	"America/Thule",
	"America/Tijuana",
	"America/Toronto",
	"America/Tortola",
	"America/Vancouver",
	"America/Virgin",
	"America/Whitehorse",
	"America/Winnipeg",
	"America/Yakutat",
	"America/Yellowknife",
	"Antarctica/Davis",
	"Antarctica/DumontDUrville",
	"Antarctica/Mawson",
	"Antarctica/McMurdo",
	"Antarctica/Palmer",
	"Antarctica/Syowa",
	"Antarctica/Troll",
	"Antarctica/Vostok",
	"Arctic/Longyearbyen",
	"Asia/Aden",
	"Asia/Almaty",
	"Asia/Amman",
	"Asia/Anadyr",
	"Asia/Aqtau",
	"Asia/Aqtobe",
	"Asia/Ashgabat",
	"Asia/Atyrau",
	"Asia/Baghdad",
	"Asia/Bahrain",
	"Asia/Baku",
	"Asia/Bangkok",
	"Asia/Barnaul",
	"Asia/Beirut",
	"Asia/Bishkek",
	"Asia/Brunei",
	"Asia/Calcutta",
	"Asia/Chita",
	"Asia/Choibalsan",
	"Asia/Chongqing",
	"Asia/Colombo",
	"Asia/Damascus",
	"Asia/Dhaka",
	"Asia/Dili",
	"Asia/Dubai",
	"Asia/Dushanbe",
	"Asia/Famagusta",
	"Asia/Gaza",
	"Asia/Hebron",
	"Asia/Ho_Chi_Minh",
	"Asia/Hong_Kong",
	"Asia/Hovd",
	"Asia/Irkutsk",
	"Asia/Istanbul",
	"Asia/Jakarta",
	"Asia/Jayapura",
	"Asia/Jerusalem",
	"Asia/Kabul",
	"Asia/Kamchatka",
	"Asia/Karachi",
	"Asia/Kathmandu",
	"Asia/Katmandu",
	"Asia/Khandyga",
	"Asia/Kolkata",
	"Asia/Krasnoyarsk",
	"Asia/Kuala_Lumpur",
	"Asia/Kuching",
	"Asia/Kuwait",
	"Asia/Macao",
	"Asia/Macau",
	"Asia/Magadan",
	"Asia/Makassar",
	"Asia/Manila",
	"Asia/Muscat",
	"Asia/Nicosia",
	"Asia/Novokuznetsk",
	"Asia/Novosibirsk",
	"Asia/Omsk",
	"Asia/Oral",
	"Asia/Phnom_Penh",
	"Asia/Pontianak",
	"Asia/Pyongyang",
	"Asia/Qatar",
	"Asia/Qostanay",
	"Asia/Qyzylorda",
	"Asia/Rangoon",
	"Asia/Riyadh",
	"Asia/Sakhalin",
	"Asia/Samarkand",
	"Asia/Seoul",
	"Asia/Shanghai",
	"Asia/Singapore",
	"Asia/Srednekolymsk",
	"Asia/Taipei",
	"Asia/Tashkent",
	"Asia/Tbilisi",
	"Asia/Tehran",
	"Asia/Thimphu",
	"Asia/Tokyo",
	"Asia/Tomsk",
	"Asia/Ujung_Pandang",
	"Asia/Ulaanbaatar",
	"Asia/Urumqi",
	"Asia/Vientiane",
	"Asia/Vladivostok",
	"Asia/Yakutsk",
	"Asia/Yangon",
	"Asia/Yekaterinburg",
	"Asia/Yerevan",
	"Atlantic/Azores",
	"Atlantic/Bermuda",
	"Atlantic/Canary",
	"Atlantic/Cape_Verde",
	"Atlantic/Faeroe",
	"Atlantic/Faroe",
	"Atlantic/Madeira",
	"Atlantic/Reykjavik",
	"Atlantic/South_Georgia",
	"Atlantic/St_Helena",
	"Atlantic/Stanley",
	"Australia/Adelaide",
	"Australia/Brisbane",
	"Australia/Broken_Hill",
	"Australia/Canberra",
	"Australia/Darwin",
	"Australia/Eucla",
	"Australia/Hobart",
	"Australia/Lindeman",
	"Australia/NSW",
	"Australia/North",
	"Australia/Perth",
	"Australia/Queensland",
	"Australia/South",
	"Australia/Sydney",
	"Australia/Tasmania",
	"Australia/Victoria",
	"Australia/West",
	"Chile/Continental",
	"Europe/Amsterdam",
	"Europe/Andorra",
	"Europe/Astrakhan",
	"Europe/Athens",
	"Europe/Belgrade",
	"Europe/Berlin",
	"Europe/Bratislava",
	"Europe/Brussels",
	"Europe/Bucharest",
	"Europe/Budapest",
	"Europe/Busingen",
	"Europe/Chisinau",
	"Europe/Copenhagen",
	"Europe/Dublin",
	"Europe/Gibraltar",
	"Europe/Guernsey",
	"Europe/Helsinki",
	"Europe/Isle_of_Man",
	"Europe/Istanbul",
	"Europe/Jersey",
	"Europe/Kaliningrad",
	"Europe/Kiev",
	"Europe/Kirov",
	"Europe/Kyiv",
	"Europe/Lisbon",
	"Europe/Ljubljana",
	"Europe/London",
	"Europe/Luxembourg",
	"Europe/Madrid",
	"Europe/Malta",
	"Europe/Mariehamn",
	"Europe/Minsk",
	"Europe/Monaco",
	"Europe/Moscow",
	"Europe/Oslo",
	"Europe/Paris",
	"Europe/Podgorica",
	"Europe/Prague",
	"Europe/Riga",
	"Europe/Rome",
	"Europe/Samara",
	"Europe/San_Marino",
	"Europe/Sarajevo",
	"Europe/Saratov",
	"Europe/Simferopol",
	"Europe/Skopje",
	"Europe/Sofia",
	"Europe/Stockholm",
	"Europe/Tallinn",
	"Europe/Tirane",
	"Europe/Ulyanovsk",
	"Europe/Vaduz",
	"Europe/Vatican",
	"Europe/Vienna",
	"Europe/Vilnius",
	"Europe/Volgograd",
	"Europe/Warsaw",
	"Europe/Zagreb",
	"Europe/Zurich",
	"Indian/Antananarivo",
	"Indian/Chagos",
	"Indian/Christmas",
	"Indian/Cocos",
	"Indian/Comoro",
	"Indian/Kerguelen",
	"Indian/Mahe",
	"Indian/Maldives",
	"Indian/Mauritius",
	"Indian/Mayotte",
	"Indian/Reunion",
	"Pacific/Apia",
	"Pacific/Auckland",
	"Pacific/Chatham",
	"Pacific/Chuuk",
	"Pacific/Easter",
	"Pacific/Efate",
	"Pacific/Enderbury",
	"Pacific/Fakaofo",
	"Pacific/Fiji",
	"Pacific/Funafuti",
	"Pacific/Galapagos",
	"Pacific/Gambier",
	"Pacific/Guadalcanal",
	"Pacific/Guam",
	"Pacific/Honolulu",
	"Pacific/Kanton",
	"Pacific/Kiritimati",
	"Pacific/Kosrae",
	"Pacific/Kwajalein",
	"Pacific/Majuro",
	"Pacific/Marquesas",
	"Pacific/Midway",
	"Pacific/Nauru",
	"Pacific/Niue",
	"Pacific/Norfolk",
	"Pacific/Noumea",
	"Pacific/Pago_Pago",
	"Pacific/Palau",
	"Pacific/Pitcairn",
	"Pacific/Pohnpei",
	"Pacific/Ponape",
	"Pacific/Port_Moresby",
	"Pacific/Rarotonga",
	"Pacific/Saipan",
	"Pacific/Samoa",
	"Pacific/Tahiti",
	"Pacific/Tarawa",
	"Pacific/Tongatapu",
	"Pacific/Truk",
	"Pacific/Wake",
	"Pacific/Wallis",
	"Pacific/Yap",
	"US/Samoa",
}

var cityTimeZoneMap = map[string]string{
	"Shanghai":    "Asia/Shanghai",
	"Los_Angeles": "America/Los_Angeles",
	"Chicago":     "America/Chicago",
	"New_York":    "America/New_York",
	"London":      "Europe/London",
	"Paris":       "Europe/Paris",
	"Moscow":      "Europe/Moscow",
	"Tokyo":       "Asia/Tokyo",
	"Singapore":   "Asia/Singapore",
}

func (tcs *TimeConvertServer) ParseCommand(stdin string) {
	tcs.Required = make(map[string]string)
	rawParts := strings.Split(stdin, " ")
	parts := utils.RemoveElements(rawParts, "").([]string)
	for i := 0; i < len(parts); i++ {
		arg := parts[i]
		if tcs.Required["sourceTime"] == "" {
			switch strings.ToUpper(arg) {
			case "--MS":
				tcs.Required["precision"] = "MilliSecond"
				continue
			case "--TZ":
				if i+1 < len(parts) {
					i++
					tcs.Options = make(map[string]string)
					ok, timeZone := validationSupportTimeZone(parts[i])
					if ok {
						tcs.Options["timeZone"] = timeZone
						continue
					}
					tcs.Options["timeZone"] = "ERR"
				}
				continue
			}
		}
		tcs.Required["sourceTime"] += arg + " "
	}
	tcs.Required["sourceTime"] = strings.TrimRight(tcs.Required["sourceTime"], " ")
	if tcs.Required["precision"] == "" {
		tcs.Required["precision"] = "Second"
	}
}

func (tcs *TimeConvertServer) ExecuteCommand(c *gin.Context) {
	timeZone := tcs.Options["timeZone"]
	if timeZone == "ERR" {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Unsupported time zones"})
		return
	}
	// todo: If no timezone information is passed in, use the geoip2 module to get the time zone based on the ip address
	location, err := time.LoadLocation(timeZone)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Unsupported time zones"})
	}

	precision := tcs.Required["precision"]
	sourceTime := tcs.Required["sourceTime"]
	var tTime *time.Time
	if sourceTime == "" {
		tTimeObj := utils.GetCurrentTime()
		tTime = &tTimeObj
	} else {
		sourceTimeObj := utils.ParseTimeString(sourceTime)
		tTimeObj := time.Date(sourceTimeObj.Year(), sourceTimeObj.Month(), sourceTimeObj.Day(), sourceTimeObj.Hour(), sourceTimeObj.Minute(), sourceTimeObj.Second(), sourceTimeObj.Nanosecond(), location)
		tTime = &tTimeObj
	}
	if tTime == nil {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Time formats that cannot be parsed."})
		return
	}
	switch precision {
	case "MilliSecond":
		c.JSON(http.StatusOK, gin.H{"response": tTime.UnixMilli()})
	case "Second":
		c.JSON(http.StatusOK, gin.H{"response": tTime.Unix()})
	}
}

func (tcs *TimestampConvertServer) ParseCommand(stdin string) {
	tcs.Required = make(map[string]string)
	rawParts := strings.Split(stdin, " ")
	parts := utils.RemoveElements(rawParts, "").([]string)
	if len(parts) == 0 {
		return
	}
	for i := 0; i < len(parts); i++ {
		arg := parts[i]
		if tcs.Required["sourceTime"] == "" {
			switch strings.ToUpper(arg) {
			case "-MS":
				tcs.Required["precision"] = "MilliSecond"
				continue
			case "-TZ":
				if i+1 < len(parts) {
					i++
					tcs.Options = make(map[string]string)
					ok, timeZone := validationSupportTimeZone(parts[i])
					if ok {
						tcs.Options["timeZone"] = timeZone
						continue
					}
					tcs.Options["timeZone"] = "ERR"
				}
				continue
			}
		}
		numStr := strings.ReplaceAll(arg, ",", "")
		tcs.Required["sourceTime"] += numStr + " "
	}
	tcs.Required["sourceTime"] = strings.TrimRight(tcs.Required["sourceTime"], " ")
	_, err := strconv.ParseInt(tcs.Required["sourceTime"], 0, 0)
	if err != nil {
		tcs.Required["sourceTime"] = "ERR"
	}
	if tcs.Required["precision"] == "" {
		tcs.Required["precision"] = "Second"
	}
}

func (tcs *TimestampConvertServer) ExecuteCommand(c *gin.Context) {
	timeZone := tcs.Options["timeZone"]
	if timeZone == "ERR" {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Unsupported time zones"})
		return
	}
	// todo: If no timezone information is passed in, use the geoip2 module to get the time zone based on the ip address
	location, err := time.LoadLocation(timeZone)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Unsupported time zones"})
	}
	precision := tcs.Required["precision"]
	sourceTime := tcs.Required["sourceTime"]
	if sourceTime == "ERR" {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Incorrect timestamp format"})
		return
	}
	var tTime string
	if precision == "Second" {
		if sourceTime == "" {
			currentTime := utils.GetCurrentTime()
			tTime = currentTime.In(location).Format("2006-01-02 15:04:05")
		} else {
			timestamp, _ := strconv.Atoi(sourceTime)
			tTime = time.Unix(int64(timestamp), 0).In(location).Format("2006-01-02 15:04:05")
		}
	} else {
		if sourceTime == "" {
			currentTime := utils.GetCurrentTime()
			tTime = currentTime.In(location).Format("2006-01-02 15:04:05.000")
		} else {
			timestamp, _ := strconv.Atoi(sourceTime)
			tTime = time.UnixMilli(int64(timestamp)).In(location).Format("2006-01-02 15:04:05.000")
		}
	}
	c.JSON(http.StatusOK, gin.H{"response": tTime})
}

func validationSupportTimeZone(tz string) (bool, string) {
	for _, timeZone := range timeZones {
		if strings.EqualFold(timeZone, tz) {
			return true, timeZone
		}
	}
	timeZone := cityTimeZoneMap[strings.Title(tz)]
	if timeZone != "" {
		return true, timeZone
	}
	return false, ""
}
