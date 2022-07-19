# Widen Exporter Script.

## Usage

### Mapping file for Webdam to Widen Migration (for Acquia)
```bash
./widen-exporter mapping --token="dummy/somerandomgeneratedtoken" --filename="yourfilename.csv"
```

### Metadata export
```bash
./widen-exporter metadata  --token="dummy/somerandomgeneratedtoken" --filename="yourfilename.csv" --query="ff:{PDF}" metadata_field_1 metadata_field_2 metadata_custom_3
```
You can get the query string from Advanced Search in Widen UI.
All metadata fields will be the machinename of metadata fields.