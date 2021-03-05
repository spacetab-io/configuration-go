package config

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestReadConfigs(t *testing.T) {
	t.Run("Success parsing common dirs and files", func(t *testing.T) {
		os.Setenv("STAGE", "dev")
		configBytes, err := ReadConfigs("./config_examples/configuration")
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		type cfg struct {
			Debug bool `yaml:"debug"`
			Log   struct {
				Level  string `yaml:"level"`
				Format string `yaml:"format"`
			} `yaml:"log"`
			Host        string `yaml:"host"`
			Port        string `yaml:"port"`
			StringValue string `yaml:"string_test"`
			BoolValue   bool   `yaml:"bool_test"`
		}

		config := &cfg{}
		err = yaml.Unmarshal(configBytes, &config)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		refConfig := &cfg{
			Debug: true,
			Log: struct {
				Level  string `yaml:"level"`
				Format string `yaml:"format"`
			}{Level: "error", Format: "text"},
			Host:        "127.0.0.1",
			Port:        "8888",
			StringValue: "",
			BoolValue:   false,
		}

		assert.EqualValues(t, refConfig, config)
	})
	t.Run("Success parsing merging many config files in default and one file in stage", func(t *testing.T) {
		os.Setenv("STAGE", "local")
		configBytes, err := ReadConfigs("./config_examples/configuration")
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		type cfg struct {
			Debug bool `yaml:"debug"`
			Redis struct {
				Hostname string `yaml:"hostname"`
				Password string `yaml:"password"`
				Database int    `yaml:"database"`
				Port     int    `yaml:"port"`
			} `yaml:"redis"`
			Log struct {
				Level  string `yaml:"level"`
				Format string `yaml:"format"`
			} `yaml:"log"`
			Host        string `yaml:"host"`
			Port        string `yaml:"port"`
			StringValue string `yaml:"string_test"`
			BoolValue   bool   `yaml:"bool_test"`
		}

		config := &cfg{}
		err = yaml.Unmarshal(configBytes, &config)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		refConfig := &cfg{
			Debug: true,
			Redis: struct {
				Hostname string `yaml:"hostname"`
				Password string `yaml:"password"`
				Database int    `yaml:"database"`
				Port     int    `yaml:"port"`
			}{Hostname: "127.1.1.1", Password: "password", Database: 123, Port: 321},
			Log: struct {
				Level  string `yaml:"level"`
				Format string `yaml:"format"`
			}{Level: "debug", Format: "русский мат"},
			Host:        "0.0.0.0",
			Port:        "6666",
			StringValue: "not a simple string",
			BoolValue:   false,
		}

		assert.EqualValues(t, refConfig, config)
	})
	t.Run("Success parsing extraordinary case", func(t *testing.T) {
		t.SkipNow()
		os.Setenv("STAGE", "prod")
		configBytes, err := ReadConfigs("./config_examples/extraordinary")
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		type ApplicationInfo struct {
			ID        int64  `json:"id"`
			Alias     string `json:"alias"`
			Name      string `json:"name"`
			About     string `json:"about"`
			Version   string `json:"version"`
			Docs      string `json:"docs"`
			Contacts  string `json:"contacts"`
			Copyright string `json:"copyright"`
		}

		type Format string
		type CallerConfig struct {
			Disabled         bool `yaml:"hide_caller"`
			CallerSkipFrames int  `yaml:"skip_frames"`
		}
		type SentryConfig struct {
			DSN    string `yaml:"dsn"`
			Enable bool   `yaml:"enable"`
		}
		type Config struct {
			Level   string        `yaml:"level"`
			Format  Format        `yaml:"format"`
			NoColor bool          `yaml:"no_color"`
			Caller  *CallerConfig `yaml:"caller"`
			Sentry  *SentryConfig `yaml:"sentry,omitempty"`
		}
		type WebServer struct {
			Host               string `yaml:"host"`
			Port               int    `yaml:"port"`
			Mode               string `yaml:"mode"`
			HasCORS            bool   `yaml:"has_cors"`
			Compress           bool   `yaml:"compress"`
			Debug              bool   `yaml:"debug"`
			ReadTimeout        int    `yaml:"read_timeout"`
			WriteTimeout       int    `yaml:"write_timeout"`
			IdleTimeout        int    `yaml:"idle_timeout"`
			ShutdownTimeout    int    `yaml:"shutdown_timeout"`
			MaxConnsPerIP      int    `yaml:"max_conn_per_ip"`
			MaxRequestsPerConn int    `yaml:"max_req_per_conn"`
		}
		type ServiceName string
		type InternalService struct {
			URL      string        `yaml:"url"`
			Version  string        `yaml:"version"`
			Timeout  time.Duration `yaml:"timeout"`
			NSQTopic string        `yaml:"nsq_topic"`

			Enable bool `yaml:"enable"`

			GzipContent bool `yaml:"gzip_content"`
			DebugEnable bool `yaml:"debug"`
		}
		type BlockingConfig struct {
			Enable   bool          `yaml:"enable"`
			Duration time.Duration `yaml:"duration"` // seconds
			Amount   int           `yaml:"amount_of_blocked_wellbores"`
		}
		type IntervalsConfig struct {
			RequestWaitTimeout    time.Duration `yaml:"request_wait_timeout"`
			RequestRepeatInterval time.Duration `yaml:"request_repeat_interval"`
			CachePeriod           time.Duration `yaml:"cache_period"`
		}
		type SensorAlias string
		type SensorAltAlias struct {
			Custom SensorAlias
			Base   SensorAlias
		}
		type SensorsAlts []SensorAltAlias // маппинг наших алиасов и uid у поставщика

		type SensorsConfig struct {
			Alts       *SensorsAlts                `yaml:"alts"`
			Correction map[SensorAlias]interface{} `yaml:"correction"`
		}
		type ResultsTime struct {
			mx sync.RWMutex
			t  map[string]time.Time
		}
		type DrillingStatusConfig struct {
			ForWells     bool `yaml:"for_wells"`
			ForWellbores bool `yaml:"for_wellbores"`
		}
		type ImportService struct {
			Blocking          BlockingConfig       `yaml:"blocking"`
			CompanyID         *uuid.UUID           `yaml:"company_id"`
			URL               string               `yaml:"url"`
			BasicAuth         *BasicAuth           `yaml:"basic_auth"`
			Intervals         IntervalsConfig      `yaml:"intervals"`
			StartTime         string               `yaml:"start_time"`
			Sensors           SensorsConfig        `yaml:"sensors"`
			DefaultActionID   int                  `yaml:"default_action_id"`
			Enable            bool                 `yaml:"enable"`
			DebugEnable       bool                 `yaml:"debug"`
			IsEmulator        bool                 `yaml:"is_emulator"`
			Historical        bool                 `yaml:"historical"`
			UseDrillingStatus DrillingStatusConfig `yaml:"use_drilling_status"`
			Concurrent        int                  `yaml:"concurrent"`

			AccessTime    *ResultsTime
			LogsStartTime *ResultsTime
			LogsEndTime   *ResultsTime
		}
		type ActionSeverity string
		type ActionRec struct {
			ID       int            `yaml:"id"`
			Name     string         `yaml:"name"`
			Severity ActionSeverity `yaml:"severity"`
		}

		type cfg struct {
			ServiceDescription ApplicationInfo                 `yaml:"info"`
			Log                Config                          `yaml:"log"`
			WebServer          WebServer                       `yaml:"webserver"`
			Service            map[ServiceName]InternalService `yaml:"service"`
			ExternalService    map[ServiceName]ImportService   `yaml:"external"`
			MessageQueue       MessageQueue                    `yaml:"message_queue"`
			Actions            map[ServiceName][]ActionRec     `yaml:"actions"`
		}

		config := cfg{}
		err = yaml.Unmarshal(configBytes, &config)
		if !assert.NoError(t, err) {
			t.FailNow()
		}
		id := uuid.MustParse("d2087717-5565-46e8-9eaa-4e36225a0ae9")
		exp := cfg{
			ServiceDescription: ApplicationInfo{
				Alias:   "import",
				Name:    "Import data service",
				Version: "v1.0.0",
			},
			Log: Config{
				Level:   "trace",
				Format:  "text",
				NoColor: false,
				Caller: &CallerConfig{
					Disabled:         false,
					CallerSkipFrames: 2,
				},
				Sentry: &SentryConfig{
					DSN:    "",
					Enable: false,
				},
			},
			WebServer: WebServer{
				Host:               "0.0.0.0",
				Port:               8080,
				Mode:               "release",
				HasCORS:            false,
				Compress:           true,
				Debug:              true,
				ReadTimeout:        60,
				WriteTimeout:       60,
				IdleTimeout:        60,
				ShutdownTimeout:    5,
				MaxConnsPerIP:      200,
				MaxRequestsPerConn: 500,
			},
			Service: map[ServiceName]InternalService{
				ServiceName("location"): {
					URL:     "http://aid-locations",
					Version: "v1",
					Timeout: 10,
				}, ServiceName("users"): {
					URL:     "http://aid-users",
					Version: "",
					Timeout: 0,
				},
			},
			ExternalService: map[ServiceName]ImportService{
				ServiceName("randomizer"): {
					Enable:      false,
					DebugEnable: false,
					IsEmulator:  true,
					Blocking: BlockingConfig{
						Enable:   true,
						Duration: 300,
						Amount:   2,
					},
					Historical:      false,
					DefaultActionID: 51,
					Intervals: IntervalsConfig{
						RequestWaitTimeout:    15,
						RequestRepeatInterval: 15,
						CachePeriod:           30,
					},
					BasicAuth:  nil,
					URL:        "http://dp-test-dp-lukoil-witsml-generator-test/generate-data",
					Concurrent: 10,
					UseDrillingStatus: DrillingStatusConfig{
						ForWells:     true,
						ForWellbores: true,
					},
					Sensors: SensorsConfig{
						Alts: &SensorsAlts{
							{
								Base:   "dver",
								Custom: "37"},
							{Base: "dtm",
								Custom: "-1"},
							{Base: "sppa",
								Custom: "18"},
							{Base: "tqa",
								Custom: "60"},
							{Base: "wob",
								Custom: "211"},
							{Base: "woba",
								Custom: "211"},
							{Base: "rpma",
								Custom: "72"},
							{Base: "mfia",
								Custom: "23"},
							{Base: "mfoa",
								Custom: "24"},
							{Base: "hkla",
								Custom: "28"},
							{Base: "dbtm",
								Custom: "35"},
							{Base: "mtia",
								Custom: "26"},
							{Base: "mtoa",
								Custom: "27"},
							{Base: "gasa",
								Custom: "32"},
							{Base: "tvt",
								Custom: "148"},
							{Base: "tva",
								Custom: "74"},
							{Base: "bpos",
								Custom: "103"},
							{
								Base:   "mdia",
								Custom: "21",
							},
							{
								Base:   "dmea",
								Custom: "36",
							},
							{
								Base:   "actc",
								Custom: "306",
							},
							{
								Base:   "ropa",
								Custom: "123",
							},
							{
								Base:   "mdoa",
								Custom: "22",
							},
						},
						Correction: map[SensorAlias]interface{}{
							"dtm":            "UTC",
							"dtm_logs_start": "UTC",
							"dtm_logs_end":   "UTC",
						},
					},
				},
				ServiceName("lukoil"): {
					Blocking: BlockingConfig{
						Enable:   false,
						Duration: 300,
						Amount:   2,
					},
					CompanyID: &id,
					URL:       "http://mskbd-sss.srv.lukoil.com/WitsmlWebSvc/witsmlsvcent.asmx",
					BasicAuth: &BasicAuth{
						Enable:   false,
						Username: "",
						Password: "",
					},
					Intervals: IntervalsConfig{
						RequestWaitTimeout:    30,
						RequestRepeatInterval: 30,
						CachePeriod:           30,
					},
					StartTime: "",
					Sensors: SensorsConfig{
						Alts: &SensorsAlts{
							{
								Base:   "dver",
								Custom: "37"},
							{Base: "dtm",
								Custom: "-1"},
							{Base: "sppa",
								Custom: "18"},
							{Base: "tqa",
								Custom: "60"},
							{Base: "wob",
								Custom: "211"},
							{Base: "woba",
								Custom: "211"},
							{Base: "rpma",
								Custom: "72"},
							{Base: "mfia",
								Custom: "23"},
							{Base: "mfoa",
								Custom: "24"},
							{Base: "hkla",
								Custom: "28"},
							{Base: "dbtm",
								Custom: "35"},
							{Base: "mtia",
								Custom: "26"},
							{Base: "mtoa",
								Custom: "27"},
							{Base: "gasa",
								Custom: "32"},
							{Base: "tvt",
								Custom: "148"},
							{Base: "tva",
								Custom: "74"},
							{Base: "bpos",
								Custom: "103"},
							{
								Base:   "mdia",
								Custom: "21",
							},
							{
								Base:   "dmea",
								Custom: "36",
							},
							{
								Base:   "actc",
								Custom: "306",
							},
							{
								Base:   "ropa",
								Custom: "123",
							},
							{
								Base:   "mdoa",
								Custom: "22",
							},
						},
						Correction: map[SensorAlias]interface{}{
							"hkla":           1.0,
							"mfia":           1.0,
							"sppa":           1.0,
							"wob":            1.0,
							"dtm":            "UTC",
							"dtm_logs_start": "UTC",
							"dtm_logs_end":   "UTC",
						},
					},
					DefaultActionID: 0,
					Enable:          true,
					DebugEnable:     false,
					IsEmulator:      false,
					Historical:      false,
					UseDrillingStatus: DrillingStatusConfig{
						ForWells:     true,
						ForWellbores: true,
					},
					Concurrent:    10,
					AccessTime:    nil,
					LogsStartTime: nil,
					LogsEndTime:   nil,
				},
			},
			MessageQueue: MessageQueue{
				Nsq: NsqQueue{
					Enable:      true,
					NsqdPort:    4150,
					LookupdPort: 4161,
					NsqdHost:    "nsq-d",
					LookupdHost: "nsq-lookup",
					LogLevel:    "debug",
				},
			},
			Actions: map[ServiceName][]ActionRec{
				ServiceName("lukoil"): {
					{ID: 0, Name: "Не определена", Severity: "danger"},
					{ID: 2, Name: "Перебуривание скважины", Severity: "good"},
					{ID: 17, Name: "Опрессовка при испытании", Severity: "warning"},
					{ID: 22, Name: "Проработка скважины в процессе бурения", Severity: "warning"},
					{ID: 23, Name: "Смена бурового инструмента", Severity: "warning"},
					{ID: 25, Name: "Опробование забойных двигателей", Severity: "warning"},
					{ID: 26, Name: "по метеорологическим условиям", Severity: "warning"},
					{ID: 27, Name: "Переостнастка талевой системы", Severity: "danger"},
					{ID: 28, Name: "Наращивание", Severity: "warning"},
					{ID: 29, Name: "Прочие", Severity: "warning"},
					{ID: 34, Name: "Ожидание электроэнергии", Severity: "danger"},
					{ID: 36, Name: "Ожидание работ по ликвидации аварий и брака", Severity: "danger"},
					{ID: 37, Name: "Определение положения инструмента", Severity: "warning"},
					{ID: 38, Name: "Ликвидация частичного поглощения бурового раствора", Severity: "warning"},
					{ID: 39, Name: "Ожидание затвердевания цемента", Severity: "danger"},
					{ID: 44, Name: "Разбуривание цементного стакана, тех. оснастки", Severity: "warning"},
					{ID: 45, Name: "Испытание колонны на герметичность", Severity: "warning"},
					{ID: 48, Name: "Ликвидация последствий аварии", Severity: "warning"},
					{ID: 51, Name: "Бурение", Severity: "good"},
					{ID: 56, Name: "по метеорологическим условиям", Severity: "danger"},
					{ID: 59, Name: "Прочие вспомогательные работы", Severity: "warning"},
					{ID: 66, Name: "Оборудование устья скважины противовыбросовым оборудованием", Severity: "danger"},
					{ID: 71, Name: "Ожидание тампонажной техники", Severity: "danger"},
					{ID: 73, Name: "Цементирование", Severity: "warning"},
					{ID: 75, Name: "Проверка противовыбросового оборудования в процессе бурения, тревога \"Выброс\"", Severity: "warning"},
					{ID: 77, Name: "Опрессовка бурового инструмента", Severity: "warning"},
					{ID: 78, Name: "Разборка оборудования", Severity: "warning"},
					{ID: 79, Name: "Замена промытого инструмента", Severity: "warning"},
					{ID: 80, Name: "Приготовление и обработка раствора", Severity: "warning"},
					{ID: 84, Name: "Подготовительные работы к заливке зон осложнений", Severity: "warning"},
					{ID: 89, Name: "Дефектоскопия бурового инструмента", Severity: "warning"},
					{ID: 94, Name: "Проверка забоя скважины после спуска колонны", Severity: "warning"},
					{ID: 109, Name: "Установка ванны", Severity: "warning"},
					{ID: 110, Name: "Сборка бурового инструмента", Severity: "warning"},
					{ID: 112, Name: "Установка цементного моста", Severity: "warning"},
					{ID: 113, Name: "Остановка скважины на приток жидкости", Severity: "warning"},
					{ID: 114, Name: "из-за отсутствия оборудования", Severity: "danger"},
					{ID: 116, Name: "СПО", Severity: "warning"},
					{ID: 118, Name: "Зарезка и бурение новым стволом до прежнего забоя", Severity: "good"},
					{ID: 119, Name: "Заливка зон осложнений", Severity: "warning"},
					{ID: 120, Name: "Ремонт оборудования", Severity: "warning"},
					{ID: 128, Name: "Устранение отказа телесистемы", Severity: "warning"},
					{ID: 129, Name: "Проработка ствола скважины перед спуском обсадной колонны", Severity: "warning"},
					{ID: 131, Name: "Торпедирование инструмента", Severity: "warning"},
					{ID: 133, Name: "СПО при расширении ствола", Severity: "warning"},
					{ID: 138, Name: "Шаблонирование ствола скважины перед спуском обсадной колонны", Severity: "warning"},
					{ID: 139, Name: "Ремонт и замена другого оборудования и инструмента", Severity: "warning"},
					{ID: 142, Name: "Подготовительные работы перед спуском обсадной колонны", Severity: "warning"},
					{ID: 144, Name: "Разборка бурового инструмента", Severity: "warning"},
					{ID: 145, Name: "Расхаживание инструмента при его затяжках", Severity: "warning"},
					{ID: 148, Name: "Отогрев оборудования", Severity: "warning"},
					{ID: 157, Name: "Прочие", Severity: "warning"},
					{ID: 158, Name: "СПО", Severity: "warning"},
					{ID: 168, Name: "Расхаживание забойных двигателей", Severity: "warning"},
					{ID: 174, Name: "Перерыв в работах по технологии работ", Severity: "danger"},
					{ID: 175, Name: "Заключительные работы после цементирования", Severity: "danger"},
					{ID: 191, Name: "Прочие", Severity: "warning"},
					{ID: 193, Name: "СПО (Вспомогательные работы)", Severity: "warning"},
					{ID: 194, Name: "Укладка бурового инструмента на стеллажах и буровой площадке", Severity: "danger"},
					{ID: 195, Name: "Расширка", Severity: "warning"},
					{ID: 196, Name: "Ожидание геофизической партии", Severity: "danger"},
					{ID: 199, Name: "Работы ловильным инструментом", Severity: "warning"},
					{ID: 203, Name: "из-за отсутствия инструмента", Severity: "danger"},
					{ID: 206, Name: "Ликвидация полного поглощения бурового раствора", Severity: "warning"},
					{ID: 207, Name: "Исследование объекта", Severity: "warning"},
					{ID: 208, Name: "Ремонт цепей", Severity: "danger"},
					{ID: 211, Name: "Работы по ликвидации неудачных заливок обсадных колонн и устранению их негерметичности", Severity: "warning"},
					{ID: 212, Name: "Замена ремней", Severity: "danger"},
					{ID: 214, Name: "Выброс бурильных труб на мостки", Severity: "danger"},
					{ID: 218, Name: "Разбуривание цементного стакана", Severity: "warning"},
					{ID: 219, Name: "Задавка скважины в связи с проявлением", Severity: "warning"},
					{ID: 222, Name: "Оборудование устья", Severity: "danger"},
					{ID: 226, Name: "Промывка ствола скважины перед спуском обсадной колонны", Severity: "warning"},
					{ID: 233, Name: "Перетяжка талевых канатов", Severity: "danger"},
					{ID: 234, Name: "Ожидание затвердевания цемента", Severity: "danger"},
					{ID: 235, Name: "Подготовительные работы к опрессовке колонны на герметичность", Severity: "warning"},
					{ID: 236, Name: "Прочие", Severity: "warning"},
					{ID: 237, Name: "Ремонтные работы", Severity: "warning"},
					{ID: 238, Name: "Промывка скважины в процессе бурения", Severity: "warning"},
					{ID: 243, Name: "ПЗР", Severity: "warning"},
					{ID: 256, Name: "Спуск хвостовика", Severity: "warning"},
					{ID: 259, Name: "Расхаживание прихваченного инструмента", Severity: "warning"},
					{ID: 262, Name: "Устранение непрохождения обсадной колонны", Severity: "warning"},
					{ID: 269, Name: "Прочие", Severity: "warning"},
					{ID: 270, Name: "Подготовительные работы для разбуривания цементного стакана в колонне", Severity: "warning"},
					{ID: 271, Name: "Ликвидация последствий осложнений", Severity: "warning"},
					{ID: 272, Name: "Сборка оборудования", Severity: "warning"},
					{ID: 275, Name: "Подготовительные работы к цементированию колонны", Severity: "danger"},
					{ID: 276, Name: "Смена долота", Severity: "warning"},
					{ID: 280, Name: "Спуск обсадной колонны", Severity: "warning"},
					{ID: 283, Name: "Сборка ловильного инструмента", Severity: "warning"},
					{ID: 285, Name: "Определение места прихвата", Severity: "warning"},
					{ID: 288, Name: "Разборка ловильного инструмента", Severity: "warning"},
					{ID: 291, Name: "ГФР", Severity: "danger"},
					{ID: 293, Name: "Прочие", Severity: "warning"},
					{ID: 295, Name: "Забойные заливки", Severity: "warning"},
					{ID: 296, Name: "Работа с телесистемой", Severity: "warning"},
					{ID: 296, Name: "Обслуживание ВСП", Severity: "danger"},
					{ID: 301, Name: "Сборка КНБК", Severity: "danger"},
					{ID: 302, Name: "Разборка КНБК", Severity: "danger"},
					{ID: 303, Name: "Опрессовка ПВО", Severity: "warning"},
					{ID: 304, Name: "Испытание в процессе бурения", Severity: "warning"},
					{ID: 306, Name: "Опрессовка обсадной колонны", Severity: "danger"},
					{ID: 307, Name: "Опрессовка цементного кольца", Severity: "danger"},
					{ID: 308, Name: "Подъем инструмента", Severity: "warning"},
					{ID: 311, Name: "Промывка", Severity: "warning"},
					{ID: 312, Name: "Спуск инструмента", Severity: "warning"},
					{ID: 313, Name: "Работа с КНБК", Severity: "warning"},
					{ID: 314, Name: "Ремонт насоса", Severity: "warning"},
					{ID: 315, Name: "Ремонт вертлюга", Severity: "warning"},
					{ID: 316, Name: "Ремонт ПКР", Severity: "warning"},
					{ID: 375, Name: "Работа с КНБК", Severity: "danger"},
					{ID: 389, Name: "Подготовительно - вспомогательные работы при испытании", Severity: "warning"},
					{ID: 390, Name: "из-за отсутствия жидкостей", Severity: "danger"},
					{ID: 391, Name: "Силовых механизмов", Severity: "warning"},
					{ID: 392, Name: "Подготовительные работы и спуск обсадной колонны или хвостовика, в т.ч.спуск-подъем бурильных труб, промывки - промежуточные и на забое скважины.", Severity: "warning"},
					{ID: 393, Name: "Смена непригодного бурового раствора после ликвидации осложнений", Severity: "warning"},
					{ID: 394, Name: "Ликвидация брака при цементировании", Severity: "warning"},
					{ID: 395, Name: "Ожидание оборудования", Severity: "danger"},
					{ID: 396, Name: "Бурение", Severity: "good"},
					{ID: 397, Name: "Работы по освобождению зацементированных НКТ", Severity: "warning"},
					{ID: 398, Name: "Демонтаж блока ГСМ", Severity: "danger"},
					{ID: 399, Name: "Ожидание работ по ликвидации аварий и брака", Severity: "warning"},
					{ID: 400, Name: "Демонтаж силового блока", Severity: "danger"},
					{ID: 401, Name: "Устранение предписаний заказчика, в т.ч.супервайзера по бурению, надзорных органов", Severity: "warning"},
					{ID: 402, Name: "Задавливание (глушение) скважины", Severity: "warning"},
					{ID: 403, Name: "Заливка колонны", Severity: "warning"},
					{ID: 404, Name: "Ремонт противовыбросового оборудования (ПВО)", Severity: "warning"},
					{ID: 405, Name: "Промывка при испытании", Severity: "warning"},
					{ID: 406, Name: "Ожидание партии ГИС", Severity: "warning"},
					{ID: 407, Name: "Прочие", Severity: "warning"},
					{ID: 408, Name: "Перфорация обсадной колонны", Severity: "danger"},
					{ID: 409, Name: "Установка грузоподъемных элементов талевой системы", Severity: "warning"},
					{ID: 410, Name: "Ремонт электрооборудования", Severity: "danger"},
					{ID: 411, Name: "из-за отсутствия цемента", Severity: "danger"},
					{ID: 412, Name: "Демонтаж модуля котельной", Severity: "danger"},
					{ID: 413, Name: "Вспомогательные работы в период разборки буровой и демонтажа оборудования", Severity: "warning"},
					{ID: 414, Name: "Межсменные простои при испытании", Severity: "danger"},
					{ID: 415, Name: "Оборудование устья обсадных колонн (за исключением монтажа/демонтажа ПВО)", Severity: "danger"},
					{ID: 416, Name: "Работы по приготовлению, обработке, утяжелению и смене бурового раствора", Severity: "warning"},
					{ID: 417, Name: "Межсменный перерыв", Severity: "danger"},
					{ID: 418, Name: "Ремонт компрессорного оборудования", Severity: "warning"},
					{ID: 419, Name: "Ликвидация прихвата бурильного инструмента", Severity: "warning"},
					{ID: 420, Name: "Ремонт буровой лебедки", Severity: "warning"},
					{ID: 421, Name: "Ожидание цемента", Severity: "danger"},
					{ID: 422, Name: "Промывка скважины в зонах обваливающихся пород", Severity: "warning"},
					{ID: 423, Name: "Ремонт устья", Severity: "danger"},
					{ID: 424, Name: "Ликвидация брака по бурильному инструменту и элементам КНБК", Severity: "warning"},
					{ID: 425, Name: "Выкачивание воды из приема", Severity: "warning"},
					{ID: 426, Name: "Исправление брака при передвижке", Severity: "warning"},
					{ID: 427, Name: "Разбуривание цем.мостов (перед зарезкой трансп.ствола)", Severity: "warning"},
					{ID: 428, Name: "Заливка устья в случае размыва или обвала", Severity: "danger"},
					{ID: 429, Name: "ПЗР", Severity: "danger"},
					{ID: 430, Name: "Ремонт системы верхнего привода", Severity: "warning"},
					{ID: 431, Name: "Цементирование", Severity: "danger"},
					{ID: 432, Name: "Вспомогательные работы в период проведения подготовительных работ", Severity: "danger"},
					{ID: 433, Name: "Подъем", Severity: "warning"},
					{ID: 434, Name: "Двигателей", Severity: "danger"},
					{ID: 435, Name: "Ликвидация брака при зарезке нового ствола", Severity: "warning"},
					{ID: 436, Name: "Работы по химобработке раствора", Severity: "warning"},
					{ID: 437, Name: "Восстановление циркуляции", Severity: "warning"},
					{ID: 438, Name: "Ремонт наземных сооружений", Severity: "danger"},
					{ID: 439, Name: "Смена талевых канатов", Severity: "danger"},
					{ID: 440, Name: "Геофизические исследования скважин (ГИС), в т.ч.испытание пластов на трубах (ИПТ)", Severity: "danger"},
					{ID: 441, Name: "Ввод в скважину инертных материалов", Severity: "warning"},
					{ID: 442, Name: "Отсутствие воды", Severity: "danger"},
					{ID: 443, Name: "Ремонт оборудования системы очистки бурового раствора", Severity: "warning"},
					{ID: 444, Name: "Ожидание рабочей силы", Severity: "danger"},
					{ID: 445, Name: "Строительство шурфа", Severity: "danger"},
					{ID: 446, Name: "Демонтаж блока противовыбросового оборудования", Severity: "danger"},
					{ID: 447, Name: "Сборка/разборка ИП", Severity: "danger"},
					{ID: 448, Name: "Прочие вспомогательные работы", Severity: "warning"},
					{ID: 449, Name: "Переостнастка талевой системы", Severity: "danger"},
					{ID: 450, Name: "Набор воды", Severity: "warning"},
					{ID: 451, Name: "Долив в процессе спуска ИП", Severity: "warning"},
					{ID: 452, Name: "Подготовка вертлюга и ведущей трубы", Severity: "danger"},
					{ID: 453, Name: "Ликвидация нефтегазопроявлений", Severity: "warning"},
					{ID: 454, Name: "Ожидание материалов, в т.ч. обсадных труб (без цемента)", Severity: "danger"},
					{ID: 455, Name: "Заполнение скважины буровым раствором", Severity: "warning"},
					{ID: 456, Name: "Механическое бурение ствола скважины", Severity: "good"},
					{ID: 457, Name: "Проверка противовыбросового оборудования в процессе бурения, тревога \"Выброс\"", Severity: "warning"},
					{ID: 458, Name: "Расширка ствола скважины", Severity: "warning"},
					{ID: 459, Name: "Подготовка противопожарного инвентаря и СИЗ персонала", Severity: "danger"},
					{ID: 460, Name: "Вспомогательные работы (сверхнормативные) при испытании", Severity: "danger"},
					{ID: 461, Name: "Ожидание тампонажной техники", Severity: "danger"},
					{ID: 462, Name: "Ликвидация аварий с геофизическими приборами", Severity: "danger"},
					{ID: 463, Name: "Работы по разборке элементов КНБК", Severity: "danger"},
					{ID: 464, Name: "Исправление траектории ствола скважины", Severity: "warning"},
					{ID: 465, Name: "Спуск промежуточной колонны для ликвидации осложнений", Severity: "warning"},
					{ID: 466, Name: "Работа с КНБК", Severity: "warning"},
					{ID: 467, Name: "Испытание пласта", Severity: "danger"},
					{ID: 468, Name: "Поломка бурильных труб и элементов бурильной колонны", Severity: "warning"},
					{ID: 469, Name: "СПО", Severity: "warning"},
					{ID: 470, Name: "из-за отсутствия обсадных труб", Severity: "danger"},
					{ID: 471, Name: "Подготовка инструмента для свинчивания-развинчивания труб", Severity: "danger"},
					{ID: 472, Name: "ГФР", Severity: "danger"},
					{ID: 473, Name: "Передвижка", Severity: "danger"},
					{ID: 474, Name: "Отсутствие пара", Severity: "danger"},
					{ID: 475, Name: "Ликвидация аварии с долотами", Severity: "warning"},
					{ID: 476, Name: "Монтаж БУ", Severity: "danger"},
					{ID: 477, Name: "Межсменный перерыв", Severity: "danger"},
					{ID: 478, Name: "из-за отсутствия оборудования", Severity: "danger"},
					{ID: 479, Name: "Очистка промывочных амбаров", Severity: "danger"},
					{ID: 480, Name: "Прочие простои", Severity: "warning"},
					{ID: 481, Name: "Ликвидация поглощений бурового раствора", Severity: "warning"},
					{ID: 482, Name: "Перебуривание скважины", Severity: "good"},
					{ID: 483, Name: "Наращивание", Severity: "warning"},
					{ID: 484, Name: "Монтаж/демонтаж ПВО и проверка его в процессе бурения", Severity: "danger"},
					{ID: 485, Name: "Ликвидация брака при ГИС, в т.ч. ИПТ", Severity: "danger"},
					{ID: 486, Name: "Простои по причине бездорожья", Severity: "danger"},
					{ID: 487, Name: "Ликвидация желобообразования", Severity: "warning"},
					{ID: 488, Name: "Ликвидация обвалов", Severity: "warning"},
					{ID: 489, Name: "Монтаж блока ГСМ", Severity: "danger"},
					{ID: 490, Name: "Отсутствие воды", Severity: "danger"},
					{ID: 491, Name: "Заполнение скважины буровым раствором", Severity: "warning"},
					{ID: 492, Name: "Испытание пластов", Severity: "danger"},
					{ID: 493, Name: "Расширка скважины в зонах обваливающихся пород", Severity: "warning"},
					{ID: 494, Name: "Ремонт оборудования для СПО", Severity: "warning"},
					{ID: 495, Name: "Переход с воды на глинистый раствор", Severity: "warning"},
					{ID: 496, Name: "Заливка, заделка шурфа", Severity: "danger"},
					{ID: 497, Name: "Ликвидация осложнений, возникших при креплении скважин", Severity: "danger"},
					{ID: 498, Name: "Ремонт оборудования талевой системы", Severity: "danger"},
					{ID: 499, Name: "Монтаж силового блока", Severity: "danger"},
					{ID: 500, Name: "Ликвидация аварий с гидравлическими забойными двигателями (ГЗД), телесистемами", Severity: "warning"},
					{ID: 501, Name: "Шаблонирование НКТ при испытании", Severity: "danger"},
					{ID: 502, Name: "Испытание (сверхнормативное) при испытании", Severity: "danger"},
					{ID: 503, Name: "Смена забойных двигателей", Severity: "warning"},
					{ID: 504, Name: "Отсутствие рабочей силы", Severity: "danger"},
					{ID: 505, Name: "Ликвидация аварии с обсадными трубами и элементами оснастки обсадных колонн", Severity: "danger"},
					{ID: 506, Name: "Восстановление циркуляции", Severity: "danger"},
					{ID: 507, Name: "Промывка скважины в процессе бурения", Severity: "warning"},
					{ID: 508, Name: "Ликвидация обвалов и осыпей горных пород", Severity: "warning"},
					{ID: 509, Name: "Ремонт силового устройства (электродвигателя или дизеля)", Severity: "danger"},
					{ID: 510, Name: "Отогрев оборудования", Severity: "danger"},
					{ID: 511, Name: "из-за отсутствия бурового раствора", Severity: "danger"},
					{ID: 512, Name: "Ремонтные работы сверх установленного лимита на данный вид работ", Severity: "warning"},
					{ID: 513, Name: "Ликвидация обрывы, расчленения бурильных труб и элементов КНБК", Severity: "warning"},
					{ID: 514, Name: "Ликвидация отказа ГЗД", Severity: "warning"},
					{ID: 515, Name: "Устранение зарезки второго ствола при проработке", Severity: "warning"},
					{ID: 516, Name: "Монтаж вышечно-лебедочного блока", Severity: "danger"},
					{ID: 517, Name: "Ликвидация пилотного ствола", Severity: "warning"},
					{ID: 518, Name: "Прихват инструмента", Severity: "warning"},
					{ID: 519, Name: "Промывка скважины после разбуривания цемента", Severity: "warning"},
					{ID: 520, Name: "из-за бездорожья", Severity: "danger"},
					{ID: 521, Name: "Исправление брака при строительстве вышки и привышечных сооружений", Severity: "danger"},
					{ID: 522, Name: "Демонтаж вышечно-лебедочного блока", Severity: "danger"},
					{ID: 523, Name: "Ликвидация аварий при цементировании скважин", Severity: "warning"},
					{ID: 524, Name: "Вспомогательные работы при ИП", Severity: "danger"},
					{ID: 525, Name: "Падение в скважину посторонних предметов", Severity: "danger"},
					{ID: 526, Name: "Ремонт вспомогательной лебедки", Severity: "warning"},
					{ID: 527, Name: "Монтаж блока дизель-электростанции", Severity: "warning"},
					{ID: 528, Name: "Электрометрические работы связанные с определением параметров осложненных интервалов", Severity: "danger"},
					{ID: 529, Name: "Промывка при испытании", Severity: "warning"},
					{ID: 530, Name: "Отсутствие рабочей силы", Severity: "danger"},
					{ID: 531, Name: "Чистка желобов", Severity: "warning"},
					{ID: 532, Name: "Отсутствие рабочей силы", Severity: "danger"},
					{ID: 533, Name: "Монтаж блока противовыбросового оборудования", Severity: "danger"},
					{ID: 534, Name: "Проработка, промывка, шаблонировка перед спуском обсадной колонны, опрессовка бурильных труб перед спуском хвостовика.", Severity: "warning"},
					{ID: 535, Name: "Ремонтные работы", Severity: "warning"},
					{ID: 536, Name: "Устранение негерметичности ПВО", Severity: "warning"},
					{ID: 537, Name: "Ремонт ротора", Severity: "danger"},
					{ID: 538, Name: "Исправление кривизны ствола скважины", Severity: "warning"},
					{ID: 539, Name: "Ликвидация нефтегазоводопроявлений", Severity: "warning"},
					{ID: 540, Name: "Демонтаж бурового оборудования", Severity: "danger"},
					{ID: 541, Name: "Исправление брака при устройстве фундаментов", Severity: "danger"},
					{ID: 542, Name: "Прочие работы по креплению скважины", Severity: "danger"},
					{ID: 543, Name: "Исправление брака при ремонтных работах", Severity: "danger"},
					{ID: 544, Name: "Перерыв в бурении основного ствола", Severity: "danger"},
					{ID: 545, Name: "Межсменный перерыв", Severity: "danger"},
					{ID: 546, Name: "из-за отсутствия химических реагентов", Severity: "danger"},
					{ID: 547, Name: "Смена машинных ключей", Severity: "danger"},
					{ID: 548, Name: "Вспомогательные работы в период проведения работ по вышкостроению и монтажу оборудования", Severity: "danger"},
					{ID: 549, Name: "Проработка скважины в процессе бурения", Severity: "warning"},
					{ID: 550, Name: "Работы по установке цементных мостов для создания исскуственного забоя перед/после спуска обсадной колонны, в т.ч. спуск-подъем бурильных труб и ОЗЦ.", Severity: "danger"},
					{ID: 551, Name: "Зарезка и бурение новым стволом до прежнего забоя", Severity: "warning"},
					{ID: 552, Name: "Долив в процессе спуска ИП", Severity: "warning"},
					{ID: 553, Name: "Испытание в процессе бурения", Severity: "warning"},
					{ID: 554, Name: "Извлечение оборванного каната", Severity: "danger"},
					{ID: 555, Name: "Ремонтно-изоляционные работы при испытании (освоении) скважины", Severity: "danger"},
					{ID: 556, Name: "по метеорологическим условиям", Severity: "danger"},
					{ID: 557, Name: "Отсутствие рабочей силы", Severity: "danger"},
					{ID: 558, Name: "Монтаж модуля котельной", Severity: "danger"},
					{ID: 559, Name: "Бурение разгрузочных (прицельных) скважин", Severity: "good"},
					{ID: 560, Name: "по метеорологическим условиям", Severity: "danger"},
					{ID: 561, Name: "Демонтаж машинно-насосного отделения", Severity: "danger"},
					{ID: 562, Name: "Ожидание работ по ликвидации брака и аварий", Severity: "danger"},
					{ID: 563, Name: "Отсутствие оборудования и инструмента", Severity: "danger"},
					{ID: 564, Name: "Ремонтные работы (сверхнормативные) при испытании", Severity: "danger"},
					{ID: 565, Name: "Аварии с забойными двигателями", Severity: "warning"},
					{ID: 566, Name: "Ожидание электроэнергии", Severity: "danger"},
					{ID: 567, Name: "Ожидание вывода на режим", Severity: "danger"},
					{ID: 568, Name: "Работы по ликвидации осложнений сверх установленного лимита на данный вид работ", Severity: "warning"},
					{ID: 569, Name: "Монтаж и подготовка вертикального силового привода", Severity: "danger"},
					{ID: 570, Name: "Спуско-подъемные операции при проработке в зонах обваливающихся пород", Severity: "warning"},
					{ID: 571, Name: "Прочие вспомогательные работы", Severity: "warning"},
					{ID: 572, Name: "Смена непригодного бурового раствора перед повторным каратажем", Severity: "warning"},
					{ID: 573, Name: "Прочие аварии", Severity: "warning"},
					{ID: 574, Name: "из-за отсутствия утяжелителей", Severity: "danger"},
					{ID: 575, Name: "Ликвидация водопроявлений", Severity: "warning"},
					{ID: 576, Name: "Ликвидация кривизны", Severity: "warning"},
					{ID: 577, Name: "Демонтаж блока дизель-электростанции", Severity: "danger"},
					{ID: 578, Name: "Снятие пакера", Severity: "warning"},
					{ID: 579, Name: "Ожидание электроэнергии", Severity: "danger"},
					{ID: 580, Name: "Холостой рейс инструмента из-за неполадок с забойным двигателем", Severity: "warning"},
					{ID: 581, Name: "Ожидание электроэнергии", Severity: "danger"},
					{ID: 582, Name: "Подготовка силового и грузоподъемного оборудования", Severity: "warning"},
					{ID: 583, Name: "Работа с КНБК", Severity: "warning"},
					{ID: 584, Name: "Другие работы, связанные с определением параметров осложненных интервалов", Severity: "warning"},
					{ID: 585, Name: "Ликвидация падения в скважину посторонних предметов", Severity: "warning"},
					{ID: 586, Name: "Монтаж бурового и силового оборудования", Severity: "danger"},
					{ID: 587, Name: "Ожидание работ по ликвидации аварий и брака", Severity: "danger"},
					{ID: 588, Name: "Подготовительные работы к передвижке", Severity: "danger"},
					{ID: 589, Name: "Строительство буровой вышки и привышечных сооружений", Severity: "danger"},
					{ID: 590, Name: "из-за бездорожья", Severity: "danger"},
					{ID: 591, Name: "Разборка буровой вышки и привышечных сооружений", Severity: "danger"},
					{ID: 592, Name: "Работы по приготовлению раствора", Severity: "warning"},
					{ID: 593, Name: "Насосного оборудования и механизмов", Severity: "danger"},
					{ID: 594, Name: "Простой по метеоусловиям", Severity: "danger"},
					{ID: 595, Name: "Буровой вышки (мачты)", Severity: "danger"},
					{ID: 596, Name: "Герметизация устья скважины глухой планшайбой", Severity: "danger"},
					{ID: 597, Name: "Геофизические работы при испытании", Severity: "danger"},
					{ID: 598, Name: "Работы по устранению нарушений колонн", Severity: "danger"},
					{ID: 599, Name: "Ремонт оборудования", Severity: "danger"},
					{ID: 600, Name: "Разбуривание оснастки, замена промывочной жидкости и опрессовка на герметичность обсадных колонн", Severity: "warning"},
					{ID: 601, Name: "Нащупывание цементного стакана", Severity: "warning"},
					{ID: 602, Name: "Установка цементных мостов для ликвидации пилотных стволов", Severity: "warning"},
					{ID: 603, Name: "Ремонт генератора", Severity: "danger"},
					{ID: 604, Name: "Внутрисменные простои при испытании", Severity: "warning"},
					{ID: 605, Name: "Работы по сборке компановки низа бурильной колонны (КНБК),  в т.ч. настройка и опробывание элементов КНБК", Severity: "danger"},
					{ID: 606, Name: "Ожидание работ по ликвидации аварий и брака", Severity: "warning"},
					{ID: 607, Name: "Исправление брака при монтаже бурового оборудования", Severity: "danger"},
					{ID: 608, Name: "Ликвидация набухания горных пород и растворения солей", Severity: "warning"},
					{ID: 609, Name: "Ожидание электроэнергии", Severity: "danger"},
					{ID: 610, Name: "Поиски прорывов изоляции при электробурении с подъемом инструмента", Severity: "danger"},
					{ID: 611, Name: "Ликвидация отказа долота", Severity: "warning"},
					{ID: 612, Name: "Поломка долот", Severity: "warning"},
					{ID: 613, Name: "Оснащение и оборудование рабочих мест персонала", Severity: "danger"},
					{ID: 614, Name: "Отсутствие оборудования и инструмента", Severity: "danger"},
					{ID: 615, Name: "Ремонтные работы в период разборки и демонтажа оборудования", Severity: "danger"},
					{ID: 616, Name: "Отсутствие оборудования и инструмента", Severity: "danger"},
					{ID: 617, Name: "Работы ловильным инструментом", Severity: "warning"},
					{ID: 618, Name: "Ликвидация осложнений, вызванных оползнями, наводнениями, ливнями и другими стихийными бедствиями", Severity: "warning"},
					{ID: 619, Name: "Холостой рейс инструмента из-за неполадок с токоподводом", Severity: "warning"},
					{ID: 620, Name: "Спуско-подъемные операции (СПО)", Severity: "warning"},
					{ID: 621, Name: "Работы по утяжелению раствора", Severity: "warning"},
					{ID: 622, Name: "Ремонтно-изоляционные работы при испытании (освоении) скважины", Severity: "warning"},
					{ID: 623, Name: "Бурение разгрузочных (прицельных) скважин", Severity: "good"},
					{ID: 624, Name: "Демонтаж БУ", Severity: "danger"},
					{ID: 625, Name: "Отсутствие воды", Severity: "danger"},
					{ID: 626, Name: "Передвижка БУ для демонтажа", Severity: "danger"},
					{ID: 627, Name: "Зарезка транспортного ствола", Severity: "warning"},
					{ID: 628, Name: "Заключительные работы по передвижке", Severity: "danger"},
					{ID: 629, Name: "Установка и регулировка КИП", Severity: "danger"},
					{ID: 630, Name: "Неудачный цементаж", Severity: "danger"},
					{ID: 631, Name: "Ликвидация отказа телесистемы", Severity: "warning"},
					{ID: 632, Name: "Ремонт манифольда высокого давления", Severity: "danger"},
					{ID: 633, Name: "Ремонт бурового насоcа", Severity: "danger"},
					{ID: 634, Name: "Подъем", Severity: "warning"},
					{ID: 635, Name: "Монтаж машинно-насосного отделения", Severity: "danger"},
					{ID: 636, Name: "Проработка скважины в зонах обваливающихся пород", Severity: "warning"},
					{ID: 637, Name: "Задавливание (глушение) скважины", Severity: "warning"},
					{ID: 638, Name: "из-за отсутствия оборудования", Severity: "danger"},
					{ID: 639, Name: "Спуск", Severity: "warning"},
					{ID: 640, Name: "Цементирование", Severity: "warning"},
					{ID: 641, Name: "Разбуривание металла", Severity: "warning"},
					{ID: 642, Name: "Подбор ловильного инструмента", Severity: "warning"},
					{ID: 643, Name: "Заливка устья в случае размыва или обвала", Severity: "warning"},
					{ID: 644, Name: "Ремонтные работы в период проведения работ по вышкостроению и монтажу оборудования", Severity: "danger"},
					{ID: 645, Name: "Подготовка желобной системы", Severity: "danger"},
					{ID: 646, Name: "Монтаж бурового шланга", Severity: "danger"},
					{ID: 647, Name: "Прочие ремонты", Severity: "danger"},
					{ID: 648, Name: "из-за бездорожья", Severity: "danger"},
					{ID: 649, Name: "из-за отсутствия оборудования", Severity: "danger"},
					{ID: 650, Name: "Отогрев оборудования", Severity: "danger"},
					{ID: 651, Name: "Спуск", Severity: "warning"},
					{ID: 652, Name: "Цементирование обсадных колонн, в т.ч. подготовительные работы перед цементажем и ожидание затвердевания цемента (ОЗЦ)", Severity: "warning"},
					{ID: 653, Name: "Отогрев оборудования", Severity: "danger"},
					{ID: 654, Name: "Подготовительные работы к испытанию объекта", Severity: "danger"},
					{ID: 655, Name: "Устройство фундамента", Severity: "danger"},
					{ID: 656, Name: "Забутовка", Severity: "warning"},
					{ID: 657, Name: "Бурение разгрузочных (прицельных) скважин", Severity: "good"},
					{ID: 658, Name: "Отсутствие пара", Severity: "danger"},
					{ID: 659, Name: "Монтаж и подготовка ротора", Severity: "danger"},
					{ID: 660, Name: "из-за бездорожья", Severity: "danger"},
					{ID: 661, Name: "Ремонтные работы в период проведения подготовительных работ", Severity: "warning"},
					{ID: 662, Name: "Устранение несоответствия параметров бурового раствора", Severity: "warning"},
					{ID: 663, Name: "Строительство термокейсов по технологии работ", Severity: "warning"},
					{ID: 664, Name: "Отсутствие пара", Severity: "danger"},
					{ID: 665, Name: "Пакеровка ИП", Severity: "warning"},
					{ID: 666, Name: "Передвижка БУ", Severity: "danger"},
					{ID: 667, Name: "Аварии с обсадными трубами", Severity: "warning"},
					{ID: 668, Name: "Вызов притока", Severity: "danger"},
					{ID: 669, Name: "Цементирование", Severity: "danger"},
					{ID: 670, Name: "Ожидание инструмента", Severity: "danger"},
					{ID: 671, Name: "Сварочные работы", Severity: "danger"},
					{ID: 672, Name: "Отогрев оборудования", Severity: "danger"},
					{ID: 673, Name: "Демонтаж фундамента", Severity: "danger"},
					{ID: 674, Name: "Ликвидация прочего брака", Severity: "danger"},
					{ID: 675, Name: "Смена тормозных колодок", Severity: "danger"},
					{ID: 676, Name: "Отогрев пневмосистемы", Severity: "danger"},
					{ID: 677, Name: "Промывка скважины в процессе бурения", Severity: "warning"},
					{ID: 678, Name: "Переостнастка талевой системы", Severity: "danger"},
					{ID: 679, Name: "Прочие", Severity: "warning"},
					{ID: 680, Name: "Извлечение заклиненной грунтоноски", Severity: "warning"},
				},
			},
		}

		assert.Equal(t, exp.Log, config.Log)
		assert.Equal(t, exp.Service, config.Service)
		assert.EqualValues(t, exp.ExternalService[ServiceName("lukoil")].URL, config.ExternalService[ServiceName("lukoil")].URL)
	})
	t.Run("Success parsing common dirs and files with different stages", func(t *testing.T) {
		os.Setenv("STAGE", "prod")
		configBytes, err := ReadConfigs("./config_examples/configuration")
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		type cfg struct {
			Debug bool `yaml:"debug"`
			Log   struct {
				Level  string `yaml:"level"`
				Format string `yaml:"format"`
			} `yaml:"log"`
			Host string `yaml:"host"`
			Port string `yaml:"port"`
		}

		config := &cfg{}
		err = yaml.Unmarshal(configBytes, &config)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		refConfig := &cfg{
			Debug: true,
			Log: struct {
				Level  string `yaml:"level"`
				Format string `yaml:"format"`
			}{Level: "error", Format: "text"},
			Host: "127.0.0.1",
			Port: "8888",
		}

		assert.EqualValues(t, refConfig, config)
	})
	t.Run("Success parsing complex dirs and files", func(t *testing.T) {
		os.Setenv("STAGE", "development")
		configBytes, err := ReadConfigs("./config_examples/configuration2")
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		type hbParams struct {
			AreaMapping map[string]string `yaml:"area_mapping"`
			Url         string            `yaml:"url"`
			Username    string            `yaml:"username"`
			Password    string            `yaml:"password"`
		}

		type cfg struct {
			HotelbookParams hbParams `yaml:"hotelbook_params"`
			Logging         string   `yaml:"logging"`
			DefaultList     []string `yaml:"default_list"`
			Databases       struct {
				Redis struct {
					Master struct {
						Username string `yaml:"username"`
						Password string `yaml:"password"`
					} `yaml:"master"`
				} `yaml:"redis"`
			} `yaml:"databases"`
		}

		config := &cfg{}
		err = yaml.Unmarshal(configBytes, &config)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		refConfig := &cfg{
			HotelbookParams: hbParams{
				AreaMapping: map[string]string{"KRK": "Krakow", "MSK": "Moscow", "CHB": "Челябинск"},
				Url:         "https://hotelbook.com/xml_endpoint",
				Username:    "TESt_USERNAME",
				Password:    "PASSWORD",
			},
			DefaultList: []string{"bar", "baz"},
			Logging:     "info",
			Databases: struct {
				Redis struct {
					Master struct {
						Username string `yaml:"username"`
						Password string `yaml:"password"`
					} `yaml:"master"`
				} `yaml:"redis"`
			}{Redis: struct {
				Master struct {
					Username string `yaml:"username"`
					Password string `yaml:"password"`
				} `yaml:"master"`
			}{Master: struct {
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			}{Username: "R_USER", Password: "R_PASS"}}},
		}

		assert.EqualValues(t, refConfig, config)
	})
	t.Run("Success parsing symlinked files and dirs", func(t *testing.T) {
		os.Setenv("STAGE", "dev")
		configBytes, err := ReadConfigs("./config_examples/symnlinkedConfigs")
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		type cfg struct {
			Debug bool `yaml:"debug"`
			Log   struct {
				Level  string `yaml:"level"`
				Format string `yaml:"format"`
			} `yaml:"log"`
			Host string `yaml:"host"`
			Port string `yaml:"port"`
		}

		config := &cfg{}
		err = yaml.Unmarshal(configBytes, &config)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		refConfig := &cfg{
			Debug: true,
			Log: struct {
				Level  string `yaml:"level"`
				Format string `yaml:"format"`
			}{Level: "error", Format: "text"},
			Host: "127.0.0.1",
			Port: "8888",
		}

		assert.EqualValues(t, refConfig, config)
	})

	if GetEnv("IN_CONTAINER", "") == "true" {
		t.Run("Success parsing symlinked files and dirs in root", func(t *testing.T) {
			os.Setenv("STAGE", "dev")
			configBytes, err := ReadConfigs("/cfgs")
			if !assert.NoError(t, err) {
				t.FailNow()
			}

			type cfg struct {
				Debug bool `yaml:"debug"`
				Log   struct {
					Level  string `yaml:"level"`
					Format string `yaml:"format"`
				} `yaml:"log"`
				Host string `yaml:"host"`
				Port string `yaml:"port"`
			}

			config := &cfg{}
			err = yaml.Unmarshal(configBytes, &config)
			if !assert.NoError(t, err) {
				t.FailNow()
			}

			refConfig := &cfg{
				Debug: true,
				Log: struct {
					Level  string `yaml:"level"`
					Format string `yaml:"format"`
				}{Level: "error", Format: "text"},
				Host: "127.0.0.1",
				Port: "8888",
			}

			assert.EqualValues(t, refConfig, config)
		})
	}

	t.Run("Fail dir not found", func(t *testing.T) {
		_, err := ReadConfigs("")
		if !assert.Error(t, err) {
			t.FailNow()
		}
	})
	t.Run("no defaults configs", func(t *testing.T) {
		_, err := ReadConfigs("./config_examples/no_defaults")
		if !assert.Error(t, err) {
			t.FailNow()
		}
	})
	t.Run("merge errors", func(t *testing.T) {
		_, err := ReadConfigs("./config_examples/merge_error")
		if !assert.Error(t, err) {
			t.FailNow()
		}
	})
}

func TestGetEnv(t *testing.T) {
	t.Run("get env key value", func(t *testing.T) {
		os.Setenv("KEY", "VALUE")
		val := GetEnv("KEY", "")
		if !assert.Equal(t, "VALUE", val) {
			t.FailNow()
		}
	})
	t.Run("get env key value fallback", func(t *testing.T) {
		os.Setenv("KEY", "VALUE")
		val := GetEnv("KEY2", "")
		if !assert.Equal(t, "", val) {
			t.FailNow()
		}
	})
}
