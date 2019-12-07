PROFILE_SWITCHER_EXE=com.github.sbelectronics.streamdeck.profileswitcher.sdPlugin\profile_switcher.exe
DEMO_EXE=com.github.sbelectronics.streamdeck.demo.sdPlugin\demo.exe
BINCLOCK_EXE=com.github.sbelectronics.streamdeck.binclock.sdPlugin\binclock.exe

#all: com.github.sbelectronics.streamdeck.profileswitcher.sdPlugin\profile_switcher.exe \
#     com.github.sbelectronics.streamdeck.demo.sdPlugin\demo.exe \
#     com.github.sbelectronics.streamdeck.binclock.sdPlugin\binclock.exe

all: $(PROFILE_SWITCHER_EXE) $(DEMO_EXE) $(BINCLOCK_EXE)

$(PROFILE_SWITCHER_EXE):
	go build -o $(PROFILE_SWITCHER_EXE) cmd/profile_switcher.go

$(DEMO_EXE):
	go build -o $(DEMO_EXE) cmd/demo.go 

$(BINCLOCK_EXE):
	go build -o $(BINCLOCK_EXE) cmd/binclock.go

clean:
	rm $(PROFILE_SWITCHER_EXE) $(DEMO_EXE) $(BINCLOCK_EXE)
