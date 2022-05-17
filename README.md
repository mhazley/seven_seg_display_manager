## Seven Segment Display Manager

A Go application to drive two HT16K33 based Adafruit 7 segment displays.

Uses `periph` for i2c.

MQTT Interface to set Numeric (+hex) strings on two displays.

Containerised for Balena. 

### MQTT INTERFACE

`mosquitto_pub -h f7f1564.local -t display/set -m '{"display_one_string": "10.12", "display_one_string": "BEEF"}'`