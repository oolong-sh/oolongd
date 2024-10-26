SetInterval(1, function()
    print("Timer event: Running every 1 seconds. Time:", os.time())
end)

OnEvent("customEvent", function(data)
    print("Custom event received:", data)
end)

return {
    run = function()
        print("Plugin is loaded and ready.")
    end,
}
