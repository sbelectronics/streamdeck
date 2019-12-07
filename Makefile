PROFILE_SWITCHER_DIR=com.github.sbelectronics.streamdeck.profileswitcher.sdPlugin
PROFILE_SWITCHER_EXE=$(PROFILE_SWITCHER_DIR)\profile_switcher.exe

DEMO_DIR=com.github.sbelectronics.streamdeck.demo.sdPlugin
DEMO_EXE=$(DEMO_DIR)\demo.exe

BINCLOCK_DIR=com.github.sbelectronics.streamdeck.binclock.sdPlugin
BINCLOCK_EXE=$(BINCLOCK_DIR)\binclock.exe

PKG_FILES=$(wildcard pkg/**/*.go)

# redefine this!
APPDATA=C:\Users\smbaker\AppData\Roaming


all: $(PROFILE_SWITCHER_EXE) $(DEMO_EXE) $(BINCLOCK_EXE)

$(PROFILE_SWITCHER_EXE): cmd/profile_switcher/*.go $(PKG_FILES)
	go build -o $(PROFILE_SWITCHER_EXE) cmd/profile_switcher/profile_switcher.go

$(DEMO_EXE): cmd/demo/*.go $(PKG_FILES)
	go build -o $(DEMO_EXE) cmd/demo/demo.go 

$(BINCLOCK_EXE): cmd/binclock/*.go $(PKG_FILES)
	go build -o $(BINCLOCK_EXE) cmd/binclock/binclock.go

install:
	xcopy /h /k /e /c /y /i $(BINCLOCK_DIR) $(APPDATA)\Elgato\StreamDeck\Plugins\$(BINCLOCK_DIR)
	xcopy /h /k /e /c /y /i $(DEMO_DIR) $(APPDATA)\Elgato\StreamDeck\Plugins\$(DEMO_DIR)
	xcopy /h /k /e /c /y /i $(PROFILE_SWITCHER_DIR) $(APPDATA)\Elgato\StreamDeck\Plugins\$(PROFILE_SWITCHER_DIR)

clean:
	rm $(PROFILE_SWITCHER_EXE) $(DEMO_EXE) $(BINCLOCK_EXE)
