-- TODO: take in template file or contents?
--
---@param date string
---@return string
function CreateDailyNote(date)
    -- TODO: name daily note (based off of current file path? -- or take in an input)
    local notePath = "daily-" .. date .. ".md"

    local f, err = io.open(notePath, "a+")
    if not f then
        print("Error opening or creating file: ", notePath, ", Error: ", err)
        return ""
    end

    -- TODO: only add if header doesn't already exist
    f:write("# Daily Note - ", date, "\n")

    f:close()

    return notePath
end
