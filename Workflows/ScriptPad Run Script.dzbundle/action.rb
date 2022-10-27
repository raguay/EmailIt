# Dropzone Action Info
# Name: ScriptPad Run Script
# Description: This will take text or a file and run a predefined script on it.
# Handles: Files, Text
# Creator: Richard Guay
# URL: http://customct.com
# Events: Clicked, Dragged
# KeyModifiers: Command, Option, Control, Shift
# SkipConfig: No
# UniqueID: 39487392048573890256
# RunsSandboxed: Yes
# Version: 1.0
# MinDropzoneVersion: 3.0
require 'net/http'
require 'json'

def dragged
  $dz.begin("Running Scripts...")
  numitems = $items.count
  scriptname = defined?( ENV['scriptname'] ) ? ENV['scriptname'] : "Upper Case"
  dragtype = defined?( ENV['dragged_type'] ) ? ENV['dragged_type'] : "text"

  #
  # Below line switches the progress display to determinate mode so we can show progress
  #
  $dz.determinate(true)
  $dz.percent(1)
  
  #
  # Index over all of the given presentations.
  #
  for index in 0 ... numitems
    #
    # Get the note text.
    #
    note = $items[index]
    if (dragtype === 'files')
      note = File.read($items[index])
    end

    #
    # Send it to ScriptPad.
    #
    $stdout.sync = true

    data = { :script => scriptname, :text => note}
    uri = URI("http://localhost:9978/api/script/run")
    http = Net::HTTP.new(uri.host, uri.port)
    req = Net::HTTP::Put.new(uri.path, 'Content-Type' => 'application/json')
    req.body = data.to_json
    res = http.request(req)
    note = JSON.parse(res.body)["text"]
    if(dragtype === 'files')
      File.write($items[index], note)
    else
      $dz.text(note)
    end
    
    #
    # Update the percentage done.
    #
    $dz.percent((index/numitems)*100)
  end

  
  #
  # The below line tells Dropzone to end with a notification 
  # center notification with the text "Copy Complete".
  #
  $dz.finish("Copy Complete")
  
  # You should always call $dz.url or $dz.text last in your script. The below $dz.text line places text on the clipboard.
  # If you don't want to place anything on the clipboard you should still call $dz.url(false)
  $dz.url(false)
end

def clicked
    #
    # Get the script name from the environment.
    #
    scriptname = defined?( ENV['scriptname'] ) ? ENV['scriptname'] : "1"

    #
    # Get a ScriptPad script name to run from the user.
    #
    config = "*.title = ScriptPad Script
    scriptname.label = What script to you want to run?
    scriptname.type = popup\n"
     
    #
    # Query the list of scripts from ScriptPad.
    #
    uri = URI("http://localhost:9978/api/scripts/list")
    res = JSON.parse(Net::HTTP.get(uri))
    
    #
    # Create the option list from the list we get from ScriptPad.
    #
    res.map{ |x| config += "scriptname.option = #{x}\n" }
    config += "    scriptname.default = #{scriptname}"
    result = $dz.pashua(config)
    scriptname = result["scriptname"]

    #
    # Set the ScriptPad script name.
    #
    $dz.save_value("scriptname", scriptname)

    #
    # Tell the user what they selected.
    #
    $dz.finish("ScriptPad Script is '#{scriptname}'")
    $dz.url(false)
end
