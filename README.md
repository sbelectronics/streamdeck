## smbaker's Streamdeck Plugins ##
Copyright Scott M Baker
http://www.smbaker.com/

| Note: All of these plugins only work with Windows. I'm a Windows (and Linux) user, not a Mac user. Sorry, Mac peoples.

This is a collection of plugins that I wrote for my Elgato Stream Deck. They are as-is without any warranty.

* binclock - binary (actually more BCD than binary, but that is the trend) clock
* http - performs HTTP GET, POST, PUT, PATCH, or DELETE operations to a given URL
* demo - a playground for testing the golang-based plugin framework
* profile_switcher - switch profiles based on active window regex match

Once installed, all of the plugins will show up in the "Custom" section of the Streamdeck software.

## binclock

A Binary Clock is a clock that displays the time in ... binary. This plugin will display a tiny binary clock inside
a button on your Streamdeck. Why would you want to do this? It's just a gadget, or a conversation piece. Use of the
plugin is relatively straightforward, just drag the button onto your button bar and it will start displaying the clock.

Windows Download: https://github.com/sbelectronics/streamdeck/blob/master/Release/com.github.sbelectronics.streamdeck.binclock.streamDeckPlugin?raw=true

## demo

This one displays a rotating box with a counter. Push the button and the color of the box will change.

## http

This plugin will send a HTTP request to a URL when the button is pushed and released. The can be useful for integrating with home automation software, websites, or any service that can be triggered from an HTTP request. 

I use mine for integration with Autohotkey, which has a Socket library that allows it to respond to HTTP requests instead of keypresses.

Windows Download: https://github.com/sbelectronics/streamdeck/blob/master/Release/com.github.sbelectronics.streamdeck.http.streamDeckPlugin?raw=true

## profile_switcher

This is a regex-based profile switcher. I did this because I wanted different profiles for different websites. The websites run as tabs in my Chrome browser. So, whenever you open a tab, Stream Deck doesn't know anything other than that you're running an Application called "Google Chrome". I made this plugin to watch for the window title of the active window, and match based on the window title. When you switch tabs in chrome, the title changes. 

| Example: The reason I wrote this was that I wanted a custom profile when visiting Tinkercad.com in Chrome, but a different profile elsewhere.

This one is hellishly complicated. It can't be released and installed, because you'll actually need to modify some of the included files yourself. You don't have to recompile it, but you do have to edit a text file and overwrite some profiles. The reason for this complexity is that the Stream Deck API only supports switching profiles to a read-only profile that is part of the plugin. You can't switch to a user-created profile, nor can you switch to a profile that was created by another plugin. So you have to make sure _this plugin_ contains your profiles.

1. Export the profiles you want to use from the Stream Deck software, save them under the names `profileswitcher1.streamDeckProfile`, `profileSwitcher2.streamDeckProfile`, etc.
2. Download this project from github if you have not already
3. Copy the contents of the directory `com.github.sbelectronics.streamdeck.profileswitcher.sdPlugin` to `%appdata%\Elgato\StreamDeck\Plugins\com.github.sbelectronics.streamdeck.profileswitcher.sdPlugin`. This installs the plugin in your StreamDeck software.
4. Go inside the directory `%appdata%\Elgato\StreamDeck\Plugins\com.github.sbelectronics.streamdeck.profileswitcher.sdPlugin`
4. Modify patterns.json to specify your regular expressions. You can add as patterns/profiles as you like.
5. Copy the profiles you exported in step #1 into the `%appdata%\Elgato\StreamDeck\Plugins\com.github.sbelectronics.streamdeck.profileswitcher.sdPlugin` directory, overwriting the profiles that are there. Now we have your custom patterns.json here and your custom profiles here.
6. Restart the Stream Deck software (if it's already running)

When you open a window that matches the Regex defined in `patterns.json`, Stream Deck will ask you to import the profile. Say yes, and then using the Stream Deck software, switch back to your default profile. After this, the profiles will change with window activations as expected.

| Note The profiles it imports are read-only, so you can't modify them after they've been updated, but you can always export another profile (create a new one as necessary), and replace the files in the `%appdata%\Elgato\StreamDeck\Plugins\com.github.sbelectronics.streamdeck.profileswitcher.sdPlugin` directory. Then, `delete` the readonly profile in Stream Deck, restart Stream Deck, and next time you visit a window that matches the regex, it'll ask you to import the profile again.

As I said, complicated, but not quite as complicated as it sounds once you get the hang of it.

| Note: ElGato developers if you ever happen across this README, please add the ability to call `switchToProfile` on profiles that don't belong to the current plugin. It would have made life way simpler.

## Building

Although I'm a Linux developer by trade, this software was written and developed in a Windows environment. I installed the following three toolsets:

* git: https://git-scm.com/download/win, version 2.24.0.2
* go: https://golang.org/dl/, version 1.12.14, 32-bit
* gcc and related: http://win-builds.org/doku.php/download_and_installation_from_windows, version 1.5.0, 32-bit (i686)

Setup your PATH variable accordingly to point to all these tools, and then running `make` should build the makefile. Expect `go` to have to spend some time downloading dependencies via `git`. `make install` (with some customization of the Makefile for your %appdata% directory) can be used to copy the plugins to the Streamdeck software plugin directory.

| Note: All of these plugins include compiled binaries in this repository, so you have no need to build them yourself unless you wish to modify them.

| Note 2: The `binclock` and `demo` plugins even have Streamdeck installers in the `Release` directory, so all you have to do is run those `.streamDeckPlugin` files and the Streamdeck software should take it from there. 
