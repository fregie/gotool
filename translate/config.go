package main

import "google.golang.org/api/androidpublisher/v3"

var (
	originLang = "en"
	targetLang = map[string]string{
		"af":     "af",
		"sq":     "sq",
		"am":     "am",
		"ar":     "ar",
		"hy-AM":  "hy",
		"az-AZ":  "az",
		"bn-BD":  "bn",
		"eu-ES":  "eu",
		"be":     "be",
		"bg":     "bg",
		"my-MM":  "my",
		"ca":     "ca",
		"zh-HK":  "zh",
		"zh-CN":  "zh",
		"zh-TW":  "zh",
		"hr":     "hr",
		"cs-CZ":  "cs",
		"da-DK":  "da",
		"nl-NL":  "nl",
		"en-IN":  "en",
		"en-SG":  "en",
		"en-ZA":  "en",
		"en-AU":  "en",
		"en-CA":  "en",
		"en-GB":  "en",
		"et":     "et",
		"fil":    "fil",
		"fi-FI":  "fi",
		"fr-CA":  "fr",
		"fr-FR":  "fr",
		"gl-ES":  "gl",
		"ka-GE":  "ka",
		"de-DE":  "de",
		"el-GR":  "el",
		"gu":     "gu",
		"iw-IL":  "iw",
		"hi-IN":  "hi",
		"hu-HU":  "hu",
		"is-IS":  "is",
		"id":     "id",
		"it-IT":  "it",
		"ja-JP":  "ja",
		"kn-IN":  "kn",
		"kk":     "kk",
		"km-KH":  "km",
		"ko-KR":  "ko",
		"ky-KG":  "ky",
		"lo-LA":  "lo",
		"lv":     "lv",
		"lt":     "lt",
		"mk-MK":  "mk",
		"ms":     "ms",
		"ms-MY":  "ms",
		"ml-IN":  "ml",
		"mr-IN":  "mr",
		"mn-MN":  "mn",
		"ne-NP":  "ne",
		"no-NO":  "no",
		"fa":     "fa",
		"fa-AE":  "fa",
		"fa-AF":  "fa",
		"pl-PL":  "pl",
		"pt-BR":  "pt",
		"pt-PT":  "pt",
		"pa":     "pa",
		"ro":     "ro",
		"ru-RU":  "ru",
		"sr":     "sr",
		"si-LK":  "si",
		"sk":     "sk",
		"sl":     "sl",
		"es-419": "es",
		"es-ES":  "es",
		"es-US":  "es",
		"sw":     "sw",
		"sv-SE":  "sv",
		"ta-IN":  "ta",
		"te-IN":  "te",
		"th":     "th",
		"tr-TR":  "tr",
		"uk":     "uk",
		"ur":     "ur",
		"vi":     "vi",
		"zu":     "zu",
	}
	originList = androidpublisher.Listing{
		Title:            "TikVPN - Free VPN Proxy Super Fast & Secure",
		ShortDescription: "TikVPN is the best privacy & security guard VPN proxy.",
		FullDescription: `Why do you need TikVPN? TikVPN is the best & fastest VPN with safe, stable, high-speed connection! Also, it is the best privacy & security guard VPN proxy. With cutting-edge encryption and our exclusively performance boosting tune, TikVPN protects your online privacy and boosts your network speed. One tap to connect, easy and stable, safe and fast. Get TikVPN now! 

		Key features:
		√ Privacy & security guard
		√ Bank-grade encryption
		√ Hide your IP and geo location 
		√ No time or bandwidth limitation
		√ No activity & connection logs
		√ One tap to connect
		√ Servers in 40+ locations
		√ Smart location
		√ Fit Wi-Fi, LTE/4G, 3G
		√ Feel free to visit TikTok
		
		👍*Enable Apps, Streaming, and Sites 
		Can not enjoy Netflix, Hulu, Facebook, Youtube, Snapchat, Whatsapp, Telegram? Try TikVPN proxy service, it enables all apps and websites for you. You can unrestrictedly enjoy all social, music, video apps or websites, wherever you are. TikVPN’s high- speed proxy service is the ever best tool for you to enable the delightful Internet without any limit!
		
		👍*Surf the Internet Anonymously 
		In the greatest online anonymity effort, TikVPN hides your real IP and geo location(WiFi or GPS based), PIN(Personal Idification Number) so that your network identity and activity cannot be tracked or analyzed by ISP or any third party(Even us). We provides neutral and independent VPN service, and haven’t ever and will never track or analyze your any log. Just help you stay completely anonymous at office, home, airport, cafe, or anywhere else.
		
		👍*Protect Privacy & Shield WiFi Hotspot
		TikVPN secures your network traffic from any hacking or tracking while you’re connected to public WiFi hotspots. Whether your password or your personal data is encrypted or in cleartext before transferring, it will be encrypted again by TikVPN. TikVPN provides strong protection to sensitive privacy. Protecting your privacy from inappropriate attacking, tracking, or leakage is our most important mission.
		
		👍*Browse Securely with Encryption
		Stay security with bank-grade encrypted traffic between your device and our servers while connected with TikVPN. It enforces the powerful protection of your connection with strong encryption. With TikVPN proxy, you are securely encrypted to browse websites, apps or anything you would like to. All your online activities will remain completely safe. 
		
		👍*Feel Free to Visit TikTok
		TikTok is the most popular short videos app all around the world wihout any doubt. However, not all the TikTok users can watch videos on TikTok due to the geo-restrictions. Please don't worry, TikVPN can help a lot with this issue. You can change Your Location and Language in TikTok, take out the SIM Card. Due to TikTok seems to use your SIM region code to decide what you see, one method you can try is to take out your SIM card firstly, and then connect a VPN to access different content from another country. TikVPN provides lots of international servers to unblock geo-restriction on TikTok.
		
		👍*Servers in 40+ Locations
		Our high-speed VPN servers located in 40+ locations all over the world, including but not limited to America, Canada, Australia, Taiwan, Hongkong, Sweden, United Kingdom, Denmark, France, Netherlands, it will expand to more and more countries later. TikVPN has fast connection with dedicated stream servers, and tremendous bandwidth, which is convenient to use.
		
		- How to Use TikVPN? 
		[1] Download it to your device. 
		[2] Follow the instructions, launch it. 
		[3] Registration, one-tap connection. 
		[4] Run your app securely and privately. 
		
		Legal:
		https://www.tikvpn.com/terms
		https://www.tikvpn.com/privacy
		
		Contact us:
		If you have any question or suggestion, please feel free to contact us:
		Email: support@tikvpn.com
		
		Or you can check our website https://www.tikvpn.com/ for more information. We’d love to hear from you.`,
		Language: "en-US",
	}
	originListShort = androidpublisher.Listing{
		Title:            "TikVPN-Fast, Secure VPN Proxy",
		ShortDescription: "TikVPN is the best privacy&security guard VPN.",
		FullDescription: `Why do you need TikVPN? TikVPN is the best & fastest VPN with safe, stable, high-speed connection! Also, it is the best privacy & security guard VPN proxy. With cutting-edge encryption and our exclusively performance boosting tune, TikVPN protects your online privacy and boosts your network speed. One tap to connect, easy and stable, safe and fast. Get TikVPN now! 

		Key features:
		√ Privacy & security guard
		√ Bank-grade encryption
		√ Hide your IP and geo location 
		√ No time or bandwidth limitation
		√ No activity & connection logs
		√ One tap to connect
		√ Servers in 40+ locations
		√ Smart location
		√ Fit Wi-Fi, LTE/4G, 3G
		
		👍*Enable Apps, Streaming, and Sites 
		Can not enjoy Netflix, Hulu, Facebook, Youtube, Snapchat, Whatsapp, Telegram? Try TikVPN proxy service, it enables all apps and websites for you. You can unrestrictedly enjoy all social, music, video apps or websites, wherever you are. TikVPN’s high- speed proxy service is the ever best tool for you to enable the delightful Internet without any limit!
		
		👍*Surf the Internet Anonymously 
		In the greatest online anonymity effort, TikVPN hides your real IP and geo location(WiFi or GPS based), PIN(Personal Idification Number) so that your network identity and activity cannot be tracked or analyzed by ISP or any third party(Even us). We provides neutral and independent VPN service, and haven’t ever and will never track or analyze your any log. Just help you stay completely anonymous at office, home, airport, cafe, or anywhere else.
		
		👍*Protect Privacy & Shield WiFi Hotspot
		TikVPN secures your network traffic from any hacking or tracking while you’re connected to public WiFi hotspots. Whether your password or your personal data is encrypted or in cleartext before transferring, it will be encrypted again by TikVPN. TikVPN provides strong protection to sensitive privacy. Protecting your privacy from inappropriate attacking, tracking, or leakage is our most important mission.
		
		👍*Browse Securely with Encryption
		Stay security with bank-grade encrypted traffic between your device and our servers while connected with TikVPN. It enforces the powerful protection of your connection with strong encryption. With TikVPN proxy, you are securely encrypted to browse websites, apps or anything you would like to. All your online activities will remain completely safe. 
		
		👍*Servers in 40+ Locations
		Our high-speed VPN servers located in 40+ locations all over the world, including but not limited to America, Canada, Australia, Taiwan, Hongkong, Sweden, United Kingdom, Denmark, France, Netherlands, it will expand to more and more countries later. TikVPN has fast connection with dedicated stream servers, and tremendous bandwidth, which is convenient to use.
		
		Legal:
		https://www.tikvpn.com/terms
		https://www.tikvpn.com/privacy
		
		Contact us:
		If you have any question or suggestion, please feel free to contact us:
		Email: support@tikvpn.com`,
		Language: "en-US",
	}
)
