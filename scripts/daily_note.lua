-- TODO: custom templates and filenames through args
local function createDailyNote()
    local date = os.date("%m-%d-%y")
    local notePath = "daily-" .. date .. ".md"

    local f, err = io.open(notePath, "r")
    if f then
        print("File " .. notePath .. " already exists.")
        return ""
    end

    f, err = io.open(notePath, "w")
    if not f then
        print(
            "Error opening or creating file: "
                .. notePath
                .. " - Error: "
                .. err
        )
        return ""
    end

    f:write("# Daily Note - ", date, "\n\n")
    f:close()

    return notePath
end

OnEvent("dailyNoteEvent", createDailyNote)

return { run = createDailyNote }
