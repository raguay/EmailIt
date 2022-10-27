#!/usr/bin/ruby

require 'net/http'
require 'json'

Encoding.default_internal = Encoding::UTF_8
Encoding.default_external = Encoding::UTF_8

input = ENV['POPCLIP_TEXT'].to_s.strip()
noteid = ENV['POPCLIP_OPTION_NOTEID'].to_s.strip()
modifier = ENV['POPCLIP_MODIFIER_FLAGS'].to_s
append = ENV['POPCLIP_OPTION_APPEND'].to_s.strip()
keycode = ENV['POPCLIP_MODIFIER_FLAGS'].to_i

if input == '-'
   input = ''
end

if keycode == 262144
   address = 'http://localhost:9697/api/message/' + URI.escape(input)
   uri = URI(address)
   res = Net::HTTP.get(uri)
elsif keycode == 524288
   address = 'http://localhost:9697/api/message/append/' + URI.escape(input)
   uri = URI(address)
   res = Net::HTTP.get(uri)
else
   address = 'http://localhost:9978/api/note/' + noteid + '/' + append
   uri = URI(address)
   req = Net::HTTP::Put.new(uri, 'Content-Type' => 'application/json')
   req.body = { note: input }.to_json
   res = Net::HTTP.start(uri.hostname, uri.port) do |http|
     http.request(req)
   end
end
