local plugins = {}
local event_handlers = {}
local timers = {}

local plugin_paths = {
    daily_note = "./scripts/daily_note.lua",
    event_plugin = "./scripts/event_plugin.lua",
}

local function loadPlugin(name, plugin_path)
    local plugin = dofile(plugin_path)
    if plugin and type(plugin.run) == "function" then
        plugins[name] = plugin
    else
        print(
            "Error loading plugin '"
                .. name
                .. "' in '"
                .. plugin_path
                .. "': Plugin must have a 'run' function"
        )
        -- os.exit() -- NOTE: should this exit on failed plugin load?
    end
end

function LoadPlugins()
    for name, path in pairs(plugin_paths) do
        loadPlugin(name, path)
    end
end

function OnEvent(event_name, handler)
    if not event_handlers[event_name] then
        event_handlers[event_name] = {}
    end
    table.insert(event_handlers[event_name], handler)
end

function SetInterval(interval, handler)
    table.insert(timers, {
        interval = interval,
        handler = handler,
        next_run = os.time() + interval,
    })
end

function TriggerEvent(event_name, ...)
    if event_handlers[event_name] then
        for _, handler in ipairs(event_handlers[event_name]) do
            handler(...)
        end
    end
end

function ProcessTimers()
    local current_time = os.time()
    for _, timer in ipairs(timers) do
        if current_time >= timer.next_run then
            timer.handler()
            timer.next_run = current_time + timer.interval
        end
    end
end

return {
    loadPlugins = LoadPlugins,
    loadPlugin = loadPlugin,
    onEvent = OnEvent,
    setInterval = SetInterval,
    triggerEvent = TriggerEvent,
    processTimers = ProcessTimers,
}
