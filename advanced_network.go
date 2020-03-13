package wireless

type rawString string

// AdvancedNetwork represents a known network
type AdvancedNetwork struct {
	Flags []string `json:"flags"`
	ID    int      `json:"id"`

	ScanSSID          *int       `json:"scan_ssid"`
	Disabled          *int       `json:"disabled"`
	Priority          *int       `json:"priority"`
	EapolFlags        *int       `json:"eapol_flags"`
	WEPTxKeyidx       *int       `json:"wep_tx_keyidx"`
	Mode              *int       `json:"mode"`
	Frequency         *int       `json:"frequency"`
	MACSecPolicy      *int       `json:"macsec_policy"`
	WpaPtkRekey       *int       `json:"wpa_ptk_rekey"`
	IDStr             *string    `json:"id_str"`
	Identity          *string    `json:"identity"`
	AnonymousIdentity *string    `json:"anonymous_identity"`
	Password          *string    `json:"password"`
	CaCert            *string    `json:"ca_cert"`
	ClientCert        *string    `json:"client_cert"`
	PrivateKey        *string    `json:"private_key"`
	PrivateKeyPasswd  *string    `json:"private_key_passwd"`
	CaCert2           *string    `json:"ca_cert2"`
	ClientCert2       *string    `json:"client_cert2"`
	PrivateKey2       *string    `json:"private_key2"`
	PrivateKeyPasswd2 *string    `json:"private_key_passwd2"`
	Phase1            *string    `json:"phase1"`
	Phase2            *string    `json:"phase2"`
	Pin               *string    `json:"pin"`
	PCSC              *string    `json:"pcsc"`
	SSID              *string    `json:"ssid"`
	PSK               *string    `json:"psk"`
	PacFile           *string    `json:"pac_file"`
	WEPKey0           *string    `json:"wep_key0"`
	WEPKey1           *string    `json:"wep_key1"`
	WEPKey2           *string    `json:"wep_key2"`
	Pairwise          *rawString `json:"pairwise"`
	KeyMgmt           *rawString `json:"key_mgmt"`
	BSSID             *rawString `json:"bssid"`
	Proto             *rawString `json:"proto"`
	EAP               *rawString `json:"eap"`
	Group             *rawString `json:"group"`
	AuthAlg           *rawString `json:"auth_alg"`
	BSSIDWhitelist    *rawString `json:"bssid_whitelist"`
	BSSIDBlacklist    *rawString `json:"bssid_blacklist"`
}

// // SetScanSSID
// func (net *AdvancedNetwork) SetScanSSID(scanSsid int) {
// 	net.ScanSSID = &scanSsid // type=int
// }

// // SetDisabled
// func (net *AdvancedNetwork) SetDisabled(disabled int) {
// 	net.Disabled = &disabled // type=int
// }

// // SetPriority
// func (net *AdvancedNetwork) SetPriority(priority int) {
// 	net.Priority = &priority // type=int
// }

// // SetEapolFlags
// func (net *AdvancedNetwork) SetEapolFlags(eapolFlags int) {
// 	net.EapolFlags = &eapolFlags // type=int
// }

// // SetWEPTxKeyidx
// func (net *AdvancedNetwork) SetWEPTxKeyidx(weptxKeyidx int) {
// 	net.WEPTxKeyidx = &weptxKeyidx // type=int
// }

// // SetMode
// func (net *AdvancedNetwork) SetMode(mode int) {
// 	net.Mode = &mode // type=int
// }

// // SetFrequency
// func (net *AdvancedNetwork) SetFrequency(frequency int) {
// 	net.Frequency = &frequency // type=int
// }

// // SetMACSecPolicy
// func (net *AdvancedNetwork) SetMACSecPolicy(macsecPolicy int) {
// 	net.MACSecPolicy = &macsecPolicy // type=int
// }

// // SetWpaPtkRekey
// func (net *AdvancedNetwork) SetWpaPtkRekey(wpaPtkRekey int) {
// 	net.WpaPtkRekey = &wpaPtkRekey // type=int
// }

// // SetIDStr
// func (net *AdvancedNetwork) SetIDStr(iDStr string) {
// 	net.IDStr = &iDStr // type=string
// }

// // SetIdentity
// func (net *AdvancedNetwork) SetIdentity(identity string) {
// 	net.Identity = &identity // type=string
// }

// // SetAnonymousIdentity
// func (net *AdvancedNetwork) SetAnonymousIdentity(anonymousIdentity string) {
// 	net.AnonymousIdentity = &anonymousIdentity // type=string
// }

// // SetPassword
// func (net *AdvancedNetwork) SetPassword(password string) {
// 	net.Password = &password // type=string
// }

// // SetCaCert
// func (net *AdvancedNetwork) SetCaCert(caCert string) {
// 	net.CaCert = &caCert // type=string
// }

// // SetClientCert
// func (net *AdvancedNetwork) SetClientCert(clientCert string) {
// 	net.ClientCert = &clientCert // type=string
// }

// // SetPrivateKey
// func (net *AdvancedNetwork) SetPrivateKey(privateKey string) {
// 	net.PrivateKey = &privateKey // type=string
// }

// // SetPrivateKeyPasswd
// func (net *AdvancedNetwork) SetPrivateKeyPasswd(privateKeyPasswd string) {
// 	net.PrivateKeyPasswd = &privateKeyPasswd // type=string
// }

// // SetCaCert2
// func (net *AdvancedNetwork) SetCaCert2(caCert2 string) {
// 	net.CaCert2 = &caCert2 // type=string
// }

// // SetClientCert2
// func (net *AdvancedNetwork) SetClientCert2(clientCert2 string) {
// 	net.ClientCert2 = &clientCert2 // type=string
// }

// // SetPrivateKey2
// func (net *AdvancedNetwork) SetPrivateKey2(privateKey2 string) {
// 	net.PrivateKey2 = &privateKey2 // type=string
// }

// // SetPrivateKeyPasswd2
// func (net *AdvancedNetwork) SetPrivateKeyPasswd2(privateKeyPasswd2 string) {
// 	net.PrivateKeyPasswd2 = &privateKeyPasswd2 // type=string
// }

// // SetPhase1
// func (net *AdvancedNetwork) SetPhase1(phase1 string) {
// 	net.Phase1 = &phase1 // type=string
// }

// // SetPhase2
// func (net *AdvancedNetwork) SetPhase2(phase2 string) {
// 	net.Phase2 = &phase2 // type=string
// }

// // SetPin
// func (net *AdvancedNetwork) SetPin(pin string) {
// 	net.Pin = &pin // type=string
// }

// // SetPCSC
// func (net *AdvancedNetwork) SetPCSC(pCSC string) {
// 	net.PCSC = &pCSC // type=string
// }

// // SetSSID
// func (net *AdvancedNetwork) SetSSID(sSID string) {
// 	net.SSID = &sSID // type=string
// }

// // SetPSK
// func (net *AdvancedNetwork) SetPSK(pSK string) {
// 	net.PSK = &pSK // type=string
// }

// // SetPacFile
// func (net *AdvancedNetwork) SetPacFile(pacFile string) {
// 	net.PacFile = &pacFile // type=string
// }

// // SetWEPKey0
// func (net *AdvancedNetwork) SetWEPKey0(wEPKey0 string) {
// 	net.WEPKey0 = &wEPKey0 // type=string
// }

// // SetWEPKey1
// func (net *AdvancedNetwork) SetWEPKey1(wEPKey1 string) {
// 	net.WEPKey1 = &wEPKey1 // type=string
// }

// // SetWEPKey2
// func (net *AdvancedNetwork) SetWEPKey2(wEPKey2 string) {
// 	net.WEPKey2 = &wEPKey2 // type=string
// }

// // SetPairwise
// func (net *AdvancedNetwork) SetPairwise(pairwise string) {
// 	v := rawString(pairwise) // type=rawString
// 	net.Pairwise = &v
// }

// // SetKeyMgmt
// func (net *AdvancedNetwork) SetKeyMgmt(keyMgmt string) {
// 	v := rawString(keyMgmt) // type=rawString
// 	net.KeyMgmt = &v
// }

// // SetBSSID
// func (net *AdvancedNetwork) SetBSSID(bSSID string) {
// 	v := rawString(bSSID) // type=rawString
// 	net.BSSID = &v
// }

// // SetProto
// func (net *AdvancedNetwork) SetProto(proto string) {
// 	v := rawString(proto) // type=rawString
// 	net.Proto = &v
// }

// // SetEAP
// func (net *AdvancedNetwork) SetEAP(eAP string) {
// 	v := rawString(eAP) // type=rawString
// 	net.EAP = &v
// }

// // SetGroup
// func (net *AdvancedNetwork) SetGroup(group string) {
// 	v := rawString(group) // type=rawString
// 	net.Group = &v
// }

// // SetAuthAlg
// func (net *AdvancedNetwork) SetAuthAlg(authAlg string) {
// 	v := rawString(authAlg) // type=rawString
// 	net.AuthAlg = &v
// }

// // SetBSSIDWhitelist
// func (net *AdvancedNetwork) SetBSSIDWhitelist(bSSIDWhitelist string) {
// 	v := rawString(bSSIDWhitelist) // type=rawString
// 	net.BSSIDWhitelist = &v
// }

// // SetBSSIDBlacklist
// func (net *AdvancedNetwork) SetBSSIDBlacklist(bSSIDBlacklist string) {
// 	v := rawString(bSSIDBlacklist) // type=rawString
// 	net.BSSIDBlacklist = &v
// }
