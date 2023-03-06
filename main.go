package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
)

const (
	resetColor   = "\033[0m"
	boldText     = "\033[1m"
	redColor     = "\033[31m"
	greenColor   = "\033[32m"
	yellowColor  = "\033[33m"
	blueColor    = "\033[34m"
	magentaColor = "\033[35m"
	cyanColor    = "\033[36m"
	whiteColor   = "\033[37m"
)

func displayHelp() {
	helpText := boldText + yellowColor + `
  ___  _____  ____  __    ____  
 / __)(  _  )(_  _)(  )  (  _ \ 
( (_-. )(_)(   )(   )(__  )(_) )
 \___/(_____) (__) (____)(____/ 
             by @WHOISshuvam <3` + resetColor + `

` + boldText + cyanColor + " Enumerate every top-level domain of any company on the fly " + resetColor + `

` + boldText + "Usage:" + resetColor + ` ` + greenColor + boldText + "gotld [Flags]" + resetColor + `

` + boldText + "Flags:" + resetColor + `
    ` + yellowColor + boldText + "-k" + resetColor + `       Keyword to search ` + redColor + `(required)` + resetColor + `
    ` + yellowColor + boldText + "-f" + resetColor + `       File containing a list of keywords ` + magentaColor + `(optional)` + resetColor + `
    ` + yellowColor + boldText + "-o" + resetColor + `       Output filename ` + redColor + `(required)` + resetColor + `
    ` + yellowColor + boldText + "-t" + resetColor + `       Number of threads to use` + cyanColor + ` (default 5)` + resetColor + `
    ` + yellowColor + boldText + "help" + resetColor + `     Display this help message `
	fmt.Println(helpText)
}

func resolveDomain(domain string, wg *sync.WaitGroup, output *os.File) {
	defer wg.Done()

	resp, err := http.Get("https://" + domain)
	if err != nil {
		resp, err = http.Get("http://" + domain)
		if err != nil {
			return
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 || resp.StatusCode == 302 || resp.StatusCode == 303 || resp.StatusCode == 401 || resp.StatusCode == 404 || resp.StatusCode == 403 || resp.StatusCode == 500 {
		fmt.Println(resp.Request.URL.Scheme+"://"+domain, resp.StatusCode)
		output.WriteString(resp.Request.URL.Scheme + "://" + resp.Request.URL.Host + "\n")
	}
}

func main() {
	var tld []string
	var finalFile, keyword, keywordFile string
	var numThreads int
	flag.StringVar(&finalFile, "o", "", "Output Filename is required.")
	flag.StringVar(&keyword, "k", "", "Keyword is required")
	flag.StringVar(&keywordFile, "f", "", "File containing list of keywords")
	flag.IntVar(&numThreads, "t", 5, "Number of threads to use")
	flag.Parse()

	if len(os.Args) == 1 || os.Args[1] == "help" || os.Args[1] == "-h" {
		displayHelp()
		return
	}

	if keywordFile != "" {
		file, err := os.Open(keywordFile)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			tld = append(tld, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}
	} else {
		tld = []string{"1inch", "aaa", "aarp", "aave", "abarth", "abb", "abbott", "abbvie", "abc", "able", "abogado", "abudhabi", "ac", "academy", "accenture", "accountant", "accountants", "aco", "active", "actor", "ad", "adac", "ads", "adult", "ae", "aeg", "aero", "aetna", "af", "afamilycompany", "afl", "africa", "ag", "agakhan", "agency", "ai", "aig", "aigo", "airbus", "airforce", "airtel", "akdn", "al", "alfaromeo", "alibaba", "alipay", "allfinanz", "allstate", "ally", "alsace", "alstom", "am", "americanexpress", "americanfamily", "amex", "amfam", "amica", "amsterdam", "an", "analytics", "android", "anquan", "anz", "ao", "aol", "apartments", "app", "apple", "aq", "aquarelle", "ar", "arab", "aragon", "aramco", "archi", "army", "arpa", "art", "arte", "as", "asda", "asia", "associates", "at", "athleta", "attorney", "au", "auction", "audi", "audible", "audio", "auspost", "author", "auto", "autos", "avax", "avianca", "aw", "aws", "ax", "axa", "az", "azure", "ba", "baby", "baidu", "bal", "banamex", "bananarepublic", "bancor", "band", "bank", "bar", "barcelona", "barclaycard", "barclays", "barefoot", "bargains", "baseball", "basketball", "bauhaus", "bayern", "bb", "bbc", "bbt", "bbva", "bcg", "bcn", "bd", "be", "beats", "beauty", "beer", "bentley", "berlin", "best", "bestbuy", "bet", "bf", "bg", "bh", "bharti", "bi", "bible", "bid", "bike", "bing", "bingo", "bio", "bit", "biz", "bj", "bl", "black", "blackfriday", "blanco", "blockbuster", "blockchain", "blog", "bloomberg", "blue", "bm", "bms", "bmw", "bn", "bnb", "bnl", "bnpparibas", "bo", "boats", "boehringer", "bofa", "bom", "bond", "boo", "book", "booking", "boots", "bosch", "bostik", "boston", "bot", "boutique", "box", "bq", "br", "bradesco", "bridgestone", "broadway", "broker", "brother", "brussels", "bs", "bt", "bts", "budapest", "bugatti", "build", "builders", "business", "buy", "buzz", "bv", "bw", "by", "bz", "bzh", "ca", "cab", "cafe", "cal", "call", "calvinklein", "cam", "camera", "camp", "cancerresearch", "canon", "capetown", "capital", "capitalone", "car", "caravan", "cards", "care", "career", "careers", "cars", "cartier", "casa", "case", "caseih", "cash", "casino", "cat", "catering", "catholic", "cba", "cbn", "cbre", "cbs", "cc", "cd", "ceb", "celo", "center", "ceo", "cern", "cf", "cfa", "cfd", "cfin", "cg", "ch", "chanel", "channel", "charity", "chase", "chat", "cheap", "chintai", "chloe", "christmas", "chrome", "chrysler", "church", "ci", "cipriani", "circle", "cisco", "citadel", "citi", "citic", "city", "cityeats", "ck", "ckb", "cl", "claims", "cleaning", "click", "clinic", "clinique", "clothing", "cloud", "club", "clubmed", "cm", "cn", "co", "coach", "codes", "coffee", "coin", "college", "cologne", "com", "combo", "comcast", "commbank", "community", "comp", "company", "compare", "computer", "comsec", "condos", "construction", "consulting", "contact", "contractors", "cooking", "cookingchannel", "cool", "coop", "corsica", "cosmos", "country", "coupon", "coupons", "courses", "cp", "cr", "credit", "creditcard", "creditunion", "cricket", "crown", "crs", "cruise", "cruises", "crypto", "csc", "cu", "cuisinella", "curve", "cv", "cw", "cx", "cy", "cymru", "cyou", "cz", "dabur", "dad", "dance", "dao", "dapp", "data", "date", "dating", "datsun", "day", "dclk", "dds", "de", "deal", "dealer", "deals", "defi", "degree", "delivery", "dell", "deloitte", "delta", "democrat", "dental", "dentist", "desi", "design", "dev", "dfn", "dhl", "diamonds", "diet", "digital", "direct", "directory", "discount", "discover", "dish", "diy", "dj", "dk", "dm", "dnp", "do", "docs", "doctor", "dodge", "dog", "doha", "domains", "doosan", "dot", "download", "drive", "dtv", "dubai", "duck", "dunlop", "duns", "dupont", "durban", "dvag", "dvr", "dz", "earth", "eat", "ec", "eco", "edeka", "edu", "education", "ee", "eg", "eh", "email", "emc", "emerck", "energy", "engineer", "engineering", "enj", "enterprises", "env", "eos", "epost", "epson", "equipment", "er", "ericsson", "erni", "es", "esq", "estate", "esurance", "et", "eth", "etisalat", "eu", "eurovision", "eus", "events", "everbank", "exchange", "expert", "exposed", "express", "extraspace", "fage", "fail", "fairwinds", "faith", "family", "fan", "fans", "farm", "farmers", "fashion", "fast", "fedex", "feedback", "ferrari", "ferrero", "fi", "fiat", "fidelity", "fido", "fil", "film", "final", "finance", "financial", "fio", "fire", "firestone", "firmdale", "fish", "fishing", "fit", "fitness", "fj", "fk", "flickr", "flights", "flir", "florist", "flow", "flowers", "flsmidth", "fly", "fm", "fo", "foo", "food", "foodnetwork", "football", "ford", "forex", "forsale", "forum", "foundation", "fox", "fr", "free", "fresenius", "frl", "frogans", "frontdoor", "frontier", "ftr", "fujitsu", "fujixerox", "fun", "fund", "furniture", "futbol", "fyi", "ga", "gal", "gallery", "gallo", "gallup", "game", "games", "gap", "garden", "gb", "gbiz", "gd", "gdn", "ge", "gea", "gent", "genting", "george", "gf", "gg", "ggee", "gh", "ghost", "gi", "gift", "gifts", "gives", "giving", "gl", "glade", "glass", "gle", "global", "globo", "gm", "gmail", "gmbh", "gmo", "gmx", "gn", "gnt", "godaddy", "gold", "goldpoint", "golf", "goo", "goodhands", "goodyear", "goog", "google", "gop", "got", "gov", "gp", "gq", "gr", "grainger", "graphics", "gratis", "green", "gripe", "grocery", "group", "gs", "gt", "gu", "guardian", "gucci", "guge", "guide", "guitars", "guru", "gw", "gy", "hair", "hamburg", "handshake", "hangout", "haus", "hbo", "hdfc", "hdfcbank", "health", "healthcare", "hedera", "help", "helsinki", "here", "hermes", "hgtv", "hiphop", "hisamitsu", "hitachi", "hiv", "hk", "hkt", "hm", "hn", "hns", "hockey", "holdings", "holiday", "homedepot", "homegoods", "homes", "homesense", "honda", "honeywell", "horse", "hospital", "host", "hosting", "hot", "hoteles", "hotels", "hotmail", "house", "how", "hr", "hsbc", "ht", "htc", "hu", "hughes", "hyatt", "hyundai", "ibm", "icbc", "ice", "icu", "id", "ie", "ieee", "ifm", "iinet", "ikano", "il", "im", "imamat", "imdb", "immo", "immobilien", "imx", "in", "industries", "infiniti", "info", "ing", "ink", "institute", "insurance", "insure", "int", "intel", "international", "intuit", "investments", "io", "iot", "iota", "ipfs", "ipiranga", "ipns", "iq", "ir", "irish", "is", "iselect", "ismaili", "ist", "istanbul", "it", "itau", "itv", "iveco", "iwc", "jaguar", "java", "jcb", "jcp", "je", "jeep", "jetzt", "jewelry", "jio", "jlc", "jll", "jm", "jmp", "jnj", "jo", "jobs", "joburg", "jot", "joy", "jp", "jpmorgan", "jprs", "juegos", "juniper", "kaufen", "kddi", "ke", "keep", "kerryhotels", "kerrylogistics", "kerryproperties", "kfh", "kg", "kh", "ki", "kia", "kim", "kinder", "kindle", "kitchen", "kiwi", "km", "kn", "koeln", "komatsu", "kosher", "kp", "kpmg", "kpn", "kr", "krd", "kred", "kuokgroup", "kusama", "kw", "ky", "kyoto", "kz", "la", "lacaixa", "ladbrokes", "lamborghini", "lamer", "lancaster", "lancia", "lancome", "land", "landrover", "lanxess", "lasalle", "lat", "latino", "latrobe", "law", "lawyer", "lb", "lc", "lds", "lease", "leclerc", "lefrak", "legal", "lego", "lexus", "lgbt", "li", "liaison", "libre", "lidl", "life", "lifeinsurance", "lifestyle", "lighting", "like", "lilly", "limited", "limo", "lincoln", "linde", "link", "lipsy", "live", "living", "lixil", "lk", "llc", "loan", "loans", "locker", "locus", "loft", "lol", "lon", "london", "lotte", "lotto", "love", "lpl", "lplfinancial", "lr", "ls", "lt", "ltd", "ltda", "lu", "lundbeck", "lupin", "luxe", "luxury", "lv", "ly", "ma", "macys", "madrid", "maif", "maison", "makeup", "man", "management", "mango", "map", "market", "marketing", "markets", "marriott", "marshalls", "maserati", "matic", "mattel", "mba", "mc", "mcd", "mcdonalds", "mckinsey", "md", "me", "med", "media", "meet", "melbourne", "meme", "memorial", "men", "menu", "meo", "merckmsd", "metlife", "mf", "mg", "mh", "miami", "microsoft", "mil", "mini", "mint", "mit", "mitsubishi", "mk", "ml", "mlb", "mls", "mm", "mma", "mn", "mo", "mobi", "mobile", "mobily", "moda", "moe", "moi", "mom", "monash", "money", "monster", "montblanc", "mopar", "mormon", "mortgage", "moscow", "moto", "motorcycles", "mov", "movie", "movistar", "mp", "mq", "mr", "ms", "msd", "mt", "mtn", "mtpc", "mtr", "mu", "museum", "music", "mutual", "mutuelle", "mv", "mw", "mx", "my", "mz", "na", "nab", "nadex", "nagoya", "name", "nationwide", "natura", "navy", "nb", "nba", "nc", "ne", "nec", "nem", "neo", "net", "netbank", "netflix", "network", "neustar", "new", "newholland", "news", "next", "nextdirect", "nexus", "nf", "nfl", "nft", "ng", "ngo", "nhk", "ni", "nico", "nike", "nikon", "ninja", "nissan", "nissay", "nl", "no", "nokia", "northwesternmutual", "norton", "now", "nowruz", "nowtv", "np", "nr", "nra", "nrw", "ntt", "nu", "nxs", "nyc", "nz", "oasis", "obi", "observer", "off", "office", "okinawa", "olayan", "olayangroup", "oldnavy", "ollo", "om", "omega", "one", "ong", "onl", "online", "onyourside", "ooo", "open", "oracle", "orange", "org", "organic", "orientexpress", "origins", "osaka", "otsuka", "ott", "ovh", "p2p", "pa", "page", "pamperedchef", "panasonic", "panerai", "paris", "pars", "partners", "parts", "party", "passagens", "pay", "payment", "pccw", "pe", "pet", "pf", "pfizer", "pg", "ph", "pharmacy", "phd", "philips", "phone", "photo", "photography", "photos", "physio", "piaget", "pics", "pictet", "pictures", "pid", "pin", "ping", "pink", "pioneer", "pizza", "pk", "pl", "place", "play", "playstation", "plumbing", "plus", "pm", "pn", "pnc", "pohl", "poker", "politie", "polkadot", "polygon", "porn", "post", "pr", "pramerica", "praxi", "press", "prime", "pro", "prod", "productions", "prof", "progressive", "promo", "properties", "property", "protection", "pru", "prudential", "ps", "pt", "pub", "pw", "pwc", "py", "qa", "qpon", "quebec", "quest", "qvc", "racing", "radio", "radix", "raid", "rar", "re", "read", "realestate", "realtor", "realty", "recipes", "red", "redstone", "redumbrella", "rehab", "reise", "reisen", "reit", "reliance", "ren", "rent", "rentals", "repair", "report", "republican", "rest", "restaurant", "review", "reviews", "rexroth", "rich", "richardli", "ricoh", "rightathome", "ril", "rio", "rip", "rmit", "ro", "robotics", "rocher", "rocks", "rodeo", "rogers", "room", "rs", "rsvp", "ru", "rugby", "ruhr", "run", "rw", "rwe", "ryukyu", "sa", "saarland", "safe", "safety", "sakura", "sale", "salon", "samsclub", "samsung", "sandvik", "sandvikcoromant", "sanofi", "sap", "sapo", "sarl", "sas", "save", "saxo", "sb", "sbi", "sbs", "sc", "sca", "scb", "schaeffler", "schmidt", "scholarships", "school", "schule", "schwarz", "science", "scjohnson", "scor", "scot", "sd", "se", "search", "seat", "secure", "security", "seek", "select", "sener", "services", "ses", "seven", "sew", "sex", "sexy", "sfr", "sg", "sh", "shangrila", "sharp", "shaw", "shell", "shia", "shift", "shiksha", "shoes", "shop", "shopping", "shouji", "show", "showtime", "shriram", "si", "sia", "silk", "sina", "singles", "site", "sj", "sk", "skale", "ski", "skin", "sky", "skype", "sl", "sling", "sm", "smart", "smile", "sn", "sncf", "so", "soccer", "social", "softbank", "software", "sohu", "sol", "solar", "solutions", "song", "sony", "sor", "soy", "space", "spiegel", "sport", "spot", "spreadbetting", "sr", "srl", "srt", "ss", "st", "stada", "staples", "star", "starhub", "statebank", "statefarm", "statoil", "stc", "stcgroup", "stockholm", "storage", "store", "stream", "studio", "study", "style", "su", "substrate", "sucks", "supplies", "supply", "support", "surf", "surgery", "suzuki", "sv", "swarm", "swatch", "swiftcover", "swiss", "sx", "sy", "sydney", "symantec", "sys", "systems", "sz", "tab", "taipei", "talk", "taobao", "target", "tatamotors", "tatar", "tattoo", "tax", "taxi", "tc", "tci", "td", "tdk", "team", "tech", "technology", "tel", "telecity", "telefonica", "temasek", "tennis", "teva", "tez", "tf", "tg", "th", "thd", "theater", "theatre", "tiaa", "ticket", "tickets", "tienda", "tiffany", "tips", "tires", "tirol", "tj", "tjmaxx", "tjx", "tk", "tkmaxx", "tl", "tm", "tmall", "tn", "to", "today", "tokyo", "tools", "top", "toray", "toshiba", "total", "tours", "town", "toyota", "toys", "tp", "tr", "trade", "trading", "training", "travel", "travelchannel", "travelers", "travelersinsurance", "trust", "trv", "tt", "tube", "tui", "tunes", "tushu", "tv", "tvs", "tw", "tz", "ua", "ubank", "ubs", "uconnect", "ug", "uk", "um", "unicom", "university", "uno", "uol", "ups", "us", "uy", "uz", "va", "vacations", "vana", "vanguard", "vc", "ve", "vegas", "ventures", "verisign", "vermögensberater", "vermögensberatung", "versicherung", "vet", "vg", "vi", "viajes", "video", "vig", "viking", "villas", "vin", "vip", "virgin", "visa", "vision", "vista", "vistaprint", "viva", "vivo", "vlaanderen", "vn", "vodka", "volkswagen", "volvo", "vote", "voting", "voto", "voyage", "vu", "vuelos", "wales", "wallet", "walmart", "walter", "wang", "wanggou", "warman", "watch", "watches", "weather", "weatherchannel", "web3", "webcam", "weber", "website", "wed", "wedding", "weibo", "weir", "wf", "whoswho", "wien", "wiki", "williamhill", "win", "windows", "wine", "winners", "wme", "wolterskluwer", "woodside", "work", "works", "world", "wow", "ws", "wtc", "wtf", "xbox", "xerox", "xfinity", "xihuan", "xin", "xlm", "xperia", "xxx", "xyz", "yachts", "yahoo", "yamaxun", "yandex", "ye", "yodobashi", "yoga", "yokohama", "you", "youtube", "yt", "yun", "za", "zappos", "zara", "zero", "zil", "zip", "zippo", "zm", "zone", "zora", "zuerich", "zw", "δοκιμή", "ελ", "бг", "бел", "дети", "ею", "испытание", "католик", "ком", "қаз", "мкд", "мон", "москва", "онлайн", "орг", "рус", "рф", "сайт", "срб", "укр", "გე", "հայ", "‏טעסט‎", "‏קום‎", "‏آزمایشی‎", "‏إختبار‎", "‏ابوظبي‎", "‏اتصالات‎", "‏ارامكو‎", "‏الاردن‎", "‏الجزائر‎", "‏السعودية‎", "‏العليان‎", "‏المغرب‎", "‏امارات‎", "‏ایران‎", "‏بارت‎", "‏بازار‎", "‏بھارت‎", "‏بيتك‎", "‏پاکستان‎", "‏ڀارت‎", "‏تونس‎", "‏سودان‎", "‏سورية‎", "‏شبكة‎", "‏عراق‎", "‏عرب‎", "‏عمان‎", "‏فلسطين‎", "‏قطر‎", "‏كاثوليك‎", "‏كوم‎", "‏مصر‎", "‏مليسيا‎", "‏موبايلي‎", "‏موريتانيا‎", "‏موقع‎", "‏همراه‎", "कॉम", "नेट", "परीक्षा", "भारत", "भारतम्", "भारोत", "संगठन", "বাংলা", "ভারত", "ভাৰত", "ਭਾਰਤ", "ભારત", "ଭାରତ", "இந்தியா", "இலங்கை", "சிங்கப்பூர்", "பரிட்சை", "భారత్", "ಭಾರತ", "ഭാരതം", "ලංකා", "คอม", "ไทย", "グーグル", "クラウド", "コム", "ストア", "セール", "テスト", "ファッション", "ポイント", "みんな", "世界", "中信", "中国", "中國", "中文网", "企业", "佛山", "信息", "健康", "八卦", "公司", "公益", "台湾", "台灣", "商城", "商店", "商标", "嘉里", "嘉里大酒店", "在线", "大众汽车", "大拿", "天主教", "娱乐", "家電", "工行", "广东", "微博", "慈善", "我爱你", "手机", "手表", "招聘", "政务", "政府", "新加坡", "新闻", "时尚", "書籍", "机构", "测试", "淡马锡", "測試", "游戏", "澳門", "点看", "珠宝", "移动", "组织机构", "网址", "网店", "网站", "网络", "联通", "诺基亚", "谷歌", "购物", "通販", "集团", "電訊盈科", "飞利浦", "食品", "餐厅", "香格里拉", "香港"}

	}

	var output *os.File
	if finalFile != "" {
		var err error
		output, err = os.Create(finalFile)
		if err != nil {
			panic(err)
		}
		defer output.Close()
	}

	var wg sync.WaitGroup
	wg.Add(len(tld))

	queue := make(chan string, len(tld))
	for _, uros := range tld {
		queue <- keyword + "." + uros
	}
	close(queue)

	for i := 0; i < numThreads; i++ {
		go func() {
			for domain := range queue {
				resolveDomain(domain, &wg, output)
			}
		}()
	}

	wg.Wait()
}
