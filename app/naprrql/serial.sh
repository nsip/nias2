#!/bin/bash

# serial.sh
# run all reporting services

time ./naprrql --itemprint
time ./naprrql --qa
time ./naprrql --report
time ./naprrql --writingextract --wordcount
time ./naprrql --xml

