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
		Title:            "Flash VPN - Fast and Unlimited Proxy",
		ShortDescription: "One tap to connect, easy and stable, safe and fast. Just Get it!",
		FullDescription: `FlashVPN is everything you need for a 100% Great VPN App. Also, it is unlimited free VPN proxy, fastest VPN, privacy VPN and security VPN for Android. Whether you are looking for security VPN on public WiFi Hotspot or fast VPN for your Netflix show, you can get all VPN functions with FlashVPN. One tap to connect, easy and stable, safe and fast. Just Get it! 

		👍 Why choose FlashVPN?
		√ Privacy & security guard
		√ Bank-grade encryption
		√ Hide your IP and geo location
		√ No time or bandwidth limitation
		√ No activity & connection logs
		√ One tap to connect
		√ 2000+servers in 40+ locations
		√ Smart location
		√ Fit Wi-Fi, LTE/4G, 3G
		
		👍 Surf the Web Securely and Privately
		Have you ever been connected to a public Wi-Fi network and wonder if anyone, somewhere, can see your online activity? It's perfectly reasonable, considering the invasion of your privacy. With FlashVPN, you can protect your information from prying and regain a certain degree of online privacy. Protect your online privacy with a lightning fast and stable connection. There are various reasons for us to use VPN in daily life. FlashVPN is a super privacy & security guard VPN proxy and has several functions to protect your personal information from hackers, save you money when shopping online, encrypt your online data, protect yourself when using PUBLIC Wi-Fi, and stop google from tracking you.
		
		👍 Protect Your Personal Information from Hackers
		FlashVPN is a super privacy & security guard VPN proxy. It protects your online privacy and boosts your network speed. One tap to connect with easy and stable, safe and fast servers. FlashVPN will provide users with 7/24 live chat support and a risk-free 30-day money-back guarantee.
		
		In the greatest online anonymity effort, FlashVPN hides your real IP and geo location (WiFi or GPS based), So that your network identity and activity cannot be tracked or analyzed by ISP or any third party. FlashVPN also secures your network traffic from any hacking or tracking while you’re connected to public WiFi hotspots. Whether your password or your personal data is encrypted or in cleartext before transferring, it will be encrypted again by FlashVPN. FlashVPN provides strong protection to sensitive privacy. FlashVPN could also offers unlimited devices support, including Windows, Mac, iOS, and Android.
		
		👍 Encrypt Your Online Data
		Your data is exposed to the ISP and even some bad guys (maybe a hacker) when connecting to the internet, which will put you at risk, and we have to encrypt the data to prevent our personal information from leaking to the bad guys,VPN (Virtual Private Network) is the best way to encrypt your data when connecting to the internet. You can fullly encrypt your data with even one click/tap, a reliable VPN, such as FlashVPN, will hide your IP address and encrypt all the data when you are connecting to the internet.
		
		👍 Protect Yourself When Using PUBLIC Wi-Fi
		If you have to connect to a public Wi-Fi, choose a Wi-Fi in well-known places such as Starbucks, the largest local Super-Market, etc. We have to say that there is no absolutely secure public Wi-Fi, however, you can choose the Wi-Fi which is likely less suspect.
		
		The most effective way to prevent your data from the risk when using public Wi-Fi is to install a reliable VPN. A reliable VPN or Virtual Private Network will encrypt your data when connecting to the internet anytime and anywhere. 
		
		We must learn to protect ourselves wherever we are when connecting to the internet, and reliable VPNs, such as FlashVPN, will help a lot with that case.
		
		In summary, FlashVPN could entirely meet your demands when you surf the Internet, and plays a crucial role in protecting your personal information from hackers, encrypting your online data, protecting yourself when using PUBLIC Wi-Fi. If you are interested in it, then download it and try its abundant functions providing with you.`,
		Language: "en-US",
	}
	originListShort = androidpublisher.Listing{
		Title:            "FlashVPN-Fast & Secure Proxy",
		ShortDescription: "One tap to connect,safe and fast!",
		FullDescription: `FlashVPN is everything you need for a 100% Great VPN App. Also, it is unlimited free VPN proxy, fastest VPN, privacy VPN and security VPN for Android. Whether you are looking for security VPN on public WiFi Hotspot or fast VPN for your Netflix show, you can get all VPN functions with FlashVPN. One tap to connect, easy and stable, safe and fast. Just Get it! 


		Why choose FlashVPN?
		√ Privacy & security guard
		√ Bank-grade encryption
		√ Hide your IP and geo location
		√ No time or bandwidth limitation
		√ No activity & connection logs
		√ One tap to connect
		√ 2000+servers in 40+ locations
		√ Smart location
		√ Fit Wi-Fi, LTE/4G, 3G
		
		Surf the Web Securely and Privately
		Have you ever been connected to a public Wi-Fi network and wonder if anyone, somewhere, can see your online activity? It's perfectly reasonable, considering the invasion of your privacy. With FlashVPN, you can protect your information from prying and regain a certain degree of online privacy. Protect your online privacy with a lightning fast and stable connection.
		
		Protect Your Personal Information from Hackers
		FlashVPN is a super privacy & security guard VPN proxy. It protects your online privacy and boosts your network speed. One tap to connect with easy and stable, safe and fast servers. FlashVPN will provide users with 7/24 live chat support and a risk-free 30-day money-back guarantee.
		
		In the greatest online anonymity effort, FlashVPN hides your real IP and geo location (WiFi or GPS based), So that your network identity and activity cannot be tracked or analyzed by ISP or any third party. 
		
		Encrypt Your Online Data
		Your data is exposed to the ISP and even some bad guys (maybe a hacker) when connecting to the internet, which will put you at risk, and we have to encrypt the data to prevent our personal information from leaking to the bad guys,VPN (Virtual Private Network) is the best way to encrypt your data when connecting to the internet. 
		
		Protect Yourself When Using PUBLIC Wi-Fi
		If you have to connect to a public Wi-Fi, choose a Wi-Fi in well-known places such as Starbucks, the largest local Super-Market, etc. We have to say that there is no absolutely secure public Wi-Fi, however, you can choose the Wi-Fi which is likely less suspect.
		
		The most effective way to prevent your data from the risk when using public Wi-Fi is to install a reliable VPN. A reliable VPN or Virtual Private Network will encrypt your data when connecting to the internet anytime and anywhere. 
		
		We must learn to protect ourselves wherever we are when connecting to the internet, and reliable VPNs, such as FlashVPN, will help a lot with that case.
		
		In summary, FlashVPN could entirely meet your demands when you surf the Internet, and plays a crucial role in protecting your personal information from hackers, encrypting your online data, protecting yourself when using PUBLIC Wi-Fi. If you are interested in it, then download it and try its abundant functions providing with you.`,
		Language: "en-US",
	}
)
