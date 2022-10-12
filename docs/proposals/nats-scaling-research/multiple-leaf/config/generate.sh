#!/bin/bash

for i in {1..20}; do
	ACCOUNT_TEMPLATE=acc$i
	cat leaf_template.config | sed "s/ACCOUNT_TEMPLATE/${ACCOUNT_TEMPLATE}/g" > leaf-acc$i.config
	sed -i "s/ITERATOR/${i}/g" leaf-acc$i.config
done
