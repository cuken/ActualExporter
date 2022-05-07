# Actual Exporter

A small utility that helps you export your data from Actual's sqlite database.

## Usage

 - Go to the [release](https://github.com/cuken/ActualExporter/releases) section and find the correct download for your Operating System
 - You will need an export from your Actual Budget, this should provide you with a db.sqlite and a metadata.json file
 - Extract the zipped ActualExporter
 - Right now the tool only support YNAB Cloud as an export type, future export options may be supported
 - Run the following command: `actualexporter export YNAB -f [PATH TO YOUR db.sqlite] -o [OUTPUT DIR FOR GENERATED BUDGET]`
   - `-f` flag needs to be supplied with the full path to your budget:
     - i.e. `c:\budget\db.sqlite` or `/home/user/db.sqlite`
   - `-o` flag needs to be supplied a directory you would like your exported budget files to be created in:
     - i.e. `c:\budget\output\` or `/home/user/budget/out`
 - The tool should export out a CSV file per Account that can be imported into YNAB via a drag in drop with the browser.
   - See [Importing CSV files to YNAB](https://docs.youneedabudget.com/article/173-file-based-import#import) for more information