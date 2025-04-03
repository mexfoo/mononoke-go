//nolint:revive,stylecheck // It has to be that way to stay original
package packets

const (
	Version200 = 0x020000
	Version967 = 0x090607
	Version963 = 0x090603
	Version920 = 0x090200
	Version811 = 0x080101
	Version740 = 0x070400
	Version520 = 0x050200
	Version410 = 0x040100
)

const (
	DateVersion200 = 200609280
	DateVersion967 = 20210128
	DateVersion920 = 201507080
	DateVersion410 = 200701120
)

const (
	ResultSuccess                                   = 0
	ResultNotExist                                  = 1
	TS_RESULT_TOO_FAR                               = 2
	TS_RESULT_NOT_OWN                               = 3
	TS_RESULT_MISC                                  = 4
	TS_RESULT_NOT_ACTABLE                           = 5
	ResultAccessDenied                              = 6
	TS_RESULT_UNKNOWN                               = 7
	TS_RESULT_DB_ERROR                              = 8
	TS_RESULT_ALREADY_EXIST                         = 9
	TS_RESULT_NOT_ENOUGH_MONEY                      = 10
	TS_RESULT_TOO_HEAVY                             = 11
	TS_RESULT_NOT_ENOUGH_JP                         = 12
	TS_RESULT_NOT_ENOUGH_LEVEL                      = 13
	TS_RESULT_NOT_ENOUGH_JOB_LEVEL                  = 14
	TS_RESULT_NOT_ENOUGH_SKILL                      = 15
	TS_RESULT_LIMIT_MAX                             = 16
	TS_RESULT_LIMIT_MIN                             = 17
	TS_RESULT_INVALID_PASSWORD                      = 18
	TS_RESULT_INVALID_TEXT                          = 19
	TS_RESULT_NOT_ENOUGH_HP                         = 20
	TS_RESULT_NOT_ENOUGH_MP                         = 21
	TS_RESULT_COOL_TIME                             = 22
	TS_RESULT_LIMIT_WEAPON                          = 23
	TS_RESULT_LIMIT_RACE                            = 24
	TS_RESULT_LIMIT_JOB                             = 25
	TS_RESULT_LIMIT_TARGET                          = 26
	TS_RESULT_NO_SKILL                              = 27
	TS_RESULT_INVALID_ARGUMENT                      = 28
	TS_RESULT_PK_LIMIT                              = 29
	TS_RESULT_NOT_ENOUGH_ENERGY                     = 31
	TS_RESULT_NOT_ENOUGH_BULLET                     = 32
	TS_RESULT_NOT_ENOUGH_EXP                        = 33
	TS_RESULT_NOT_ENOUGH_ITEM                       = 34
	TS_RESULT_LIMIT_RIDING                          = 35
	TS_RESULT_NOT_ENOUGH_SP                         = 36
	TS_RESULT_ALREADY_STAMINA_SAVED                 = 37
	ResultTooYoung                                  = 38
	TS_RESULT_WITHDRAW_WAITING                      = 39
	TS_RESULT_REALNAME_REQUIRED                     = 40
	TS_RESULT_GAMETIME_TIRED_STAMINA_SAVER          = 41
	TS_RESULT_GAMETIME_HARMFUL_STAMINA_SAVER        = 42
	TS_RESULT_NOT_ACTABLE_IN_SIEGE_OR_RAID          = 44
	TS_RESULT_NOT_ACTABLE_IN_SECROUTE               = 45
	TS_RESULT_NOT_ACTABLE_IN_EVENTMAP               = 46
	TS_RESULT_TARGET_IN_SIEGE_OR_RAID               = 47
	TS_RESULT_TARGET_IN_SECROUTE                    = 48
	TS_RESULT_TARGET_IN_EVENTMAP                    = 49
	TS_RESULT_TOO_CHEAP                             = 50
	TS_RESULT_NOT_ACTABLE_WHILE_USING_STORAGE       = 51
	TS_RESULT_NOT_ACTABLE_WHILE_TRADING             = 52
	TS_RESULT_TOO_MUCH_MONEY                        = 53
	TS_RESULT_PASSWORD_MISMATCH                     = 54
	TS_RESULT_NOT_ACTABLE_WHILE_USING_BOOTH         = 55
	TS_RESULT_NOT_ACTABLE_IN_HUNTAHOLIC             = 56
	TS_RESULT_TARGET_IN_HUNTAHOLIC                  = 57
	TS_RESULT_NOT_ENOUGH_HUNTAHOLIC_POINT           = 58
	TS_RESULT_ACTABLE_IN_ONLY_HUNTAHOLIC            = 59
	TS_RESULT_IP_BLOCKED                            = 60
	TS_RESULT_ALREADY_IN_COMPETE                    = 61
	TS_RESULT_NOT_IN_COMPETE                        = 62
	TS_RESULT_WAITING_COMPETE_REQUEST_ANSWER        = 63
	TS_RESULT_NOT_IN_COMPETIBLE_PLACE               = 64
	TS_RESULT_TARGET_ALREADY_IN_COMPETE             = 65
	TS_RESULT_TARGET_NOT_IN_COMPETE                 = 66
	TS_RESULT_TARGET_WAITING_COMPETE_REQUEST_ANSWER = 67
	TS_RESULT_TARGET_NOT_IN_COMPETIBLE_PLACE        = 68
	TS_RESULT_NOT_ACTABLE_HERE                      = 69
	TS_RESULT_GAMETIME_LIMITED                      = 71
	TS_RESULT_NOT_ACTABLE_IN_DEATHMATCH             = 72
	TS_RESULT_ACTABLE_IN_ONLY_DEATHMATCH            = 73
	TS_RESULT_BLOCK_CHAT                            = 74
	TS_RESULT_ENHANCE_LIMIT                         = 76
	TS_RESULT_PENDING                               = 77
	TS_RESULT_NOT_ACTABLE_IN_SECRET_DUNGEON         = 78
	TS_RESULT_TARGET_IN_SECRET_DUNGEON              = 79
	TS_RESULT_ALREADY_SUPER_SAVER                   = 80
	TS_RESULT_GAMETIME_TIRED_SUPER_SAVER            = 81
	TS_RESULT_GAMETIME_HARMFUL_SUPER_SAVER          = 82
	TS_RESULT_NOT_ENOUGH_TP                         = 83
	TS_RESULT_NOT_ACTABLE_IN_INSTANCE_DUNGEON       = 84
	TS_RESULT_ACTABLE_IN_ONLY_INSTANCE_DUNGEON      = 85
	TS_RESULT_TARGET_IN_INSTANCE_DUNGEON            = 86
	TS_RESULT_TARGET_IN_DEATHMATCH                  = 87
	TS_RESULT_TARGET_IS_USING_STORAGE               = 88
	TS_RESULT_NOT_ENOUGH_AGE_PERIOD                 = 89
	TS_RESULT_ALREADY_TAMING                        = 70
	TS_RESULT_NOT_TAMABLE                           = 90
	TS_RESULT_TARGET_ALREADY_BEING_TAMED            = 91
	TS_RESULT_NOT_ENOUGH_TARGET_HP                  = 92
	TS_RESULT_NOT_ENOUGH_SUMMON_CARD                = 93
	TS_RESULT_NOT_ENOUGH_SOUL_TAMING_CARD           = 94
	TS_RESULT_NOT_ACTABLE_IN_BATTLE_ARENA           = 95
	TS_RESULT_NOT_READY                             = 96
	TS_RESULT_TARGET_IN_BATTLE_ARENA                = 97
	TS_RESULT_NOT_ACTABLE_ON_STAND_UP               = 98
	TS_RESULT_NOT_ENOUGH_ARENA_POINT                = 99
	TS_RESULT_SUCCESS_WITHOUT_NOTICE                = 101
	TS_RESULT_WEBZEN_DUPLICATE_ACCOUNT              = 102
	TS_RESULT_WEBZEN_NEED_ACCEPT_EULA               = 103
)

type Message struct {
	HeaderMessageSize     uint32
	HeaderMessageId       uint16
	HeaderMessageChecksum uint8
}

func (m *Message) GetHeaderChecksum() uint8 {
	value := uint32(0)

	value += m.HeaderMessageSize & 0xFF
	value += (m.HeaderMessageSize >> 8) & 0xFF
	value += (m.HeaderMessageSize >> 16) & 0xFF
	value += (m.HeaderMessageSize >> 24) & 0xFF

	value += uint32(m.HeaderMessageId) & 0xFF
	value += (uint32(m.HeaderMessageId) >> 8) & 0xFF

	return uint8(value) //nolint:gosec // This is fine.
}

func SetHeaderChecksum(size, id int) uint8 {
	value := 0

	value += size & 0xFF
	value += (size >> 8) & 0xFF
	value += (size >> 16) & 0xFF
	value += (size >> 24) & 0xFF

	value += id & 0xFF
	value += (id >> 8) & 0xFF

	return uint8(value) //nolint:gosec // This is fine.
}

/*
200609280
200701120
201507080
20210128
*/
