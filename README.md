# Ever Parse

Script for parsing data from Evercore Heroes after it's been exported.

## Extract UE5 mappings so we can view assets

1. Download [Process Hacker](https://processhacker.sourceforge.io/)
2. Download [UnrealMappinsDumper](https://github.com/OutTheShade/UnrealMappingsDumper)
3. Start Process Hacker
4. Start Evercore Heroes (not just the launcher)
5. Inject the DLL into the Evercore Heroes process
   1. Right click on the ProjectV-Win64-Shipping.exe process
   2. Select Miscellaneous --> Inject DLL...
   3. Select the UnrealMappingsDumper.dll file
   4. It should open a small terminal and state that the mapping file was created
6. The Mapping file is located in: `\evercore-heroes\live\ProjectV\Binaries\Win64\Mappings.usmap`

## View Evercore Heroes assets

1. Download [FModel](https://fmodel.app/)
2. Run FModel
3. Select to the evercore-heroes folder in the directory selector
4. Open Settings in the top bar
   1. Set the UE Versions drop down to UE5_1 (this might vary depending on what files you want to extract)
   2. Check the Local Mapping option and select the mappings file you created in the previous step.
   3. Navigate to the Models tab and check Auto-Export without Previewing for all the options there.
   4. Press OK
5. Open the .utoc archive and browse away! If you get an mapping error try restarting FModel

## Extract Data

1. Open the .utoc archive in FModel
2. Find the folders that are of interest and right click on them. Select to Extract as JSON and then as PNG
3. Some exports might fail so try going deeper into the folders or changing UE version.
4. A folder called Output will be located in the same place as FModel.exe and the export can be found under Exports

## Parsing the Extracted Data

1. Create a folder called "Game" in the root directory
2. Copy folders in the exported content folder `Output\Exports\ProjectV\Content` to the Game folder in root
3. Run "main .go" file
4. A directory called "parsedData" and one called "icons should be exported"
