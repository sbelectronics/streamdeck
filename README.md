## smbaker's Streamdeck Plugins ##
Copyright Scott M Baker
http://www.smbaker.com/

| Note: All of these plugins only work with Windows. I'm a Windows (and Linux) user, not a Mac user. Sorry, Mac peoples.

This is a collection of plugins that I wrote for my Elgato Stream Deck. They are as-is without any warranty.

* binclock - binary (actually more BCD than binary, but that is the trend) clock
* demo - a playground for testing the golang-based plugin framework
* profile_switcher - switch profiles based on active window regex match

## binclock

This one is relatively straightforward, just drag it onto your button bar and it will start displaying the clock.

## demo

This one displays a rotating box with a counter. Push the button and the color of the box will change.

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
