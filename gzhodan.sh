#!/bin/bash
for i in $(cat urls.txt); do 
	firefox --new-tab $i; 
done
