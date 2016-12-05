#!/bin/bash

TGS="tg1 tg2 tg3 tg4"
USERS="user1 user2 user3 user4"
LISTS="warroom"

CLI=cli/tcli.go

echo "* Creating users..."
for US in ${USERS}
do
	echo "- ${US}"
	$CLI user new ${US} ${US}@example.org
	$CLI user set fullname ${US} ${US}
	$CLI user set telephone ${US} "+15552352352"
	$CLI user set sms ${US} "+15552352352"
	$CLI user set postal ${US} "Some Street, The City"
	$CLI user set biography ${US} "So much to tell"
done

echo "* Creating TrustGroups..."
for TG in ${TGS};
do
	echo "- ${TG}"
	$CLI tg add ${TG}
	$CLI tg set groupdesc "Generated: ${TG}"

	for US in ${USERS}
	do
		echo " - adding user ${US}"
		$CLI tg member add ${TG} ${US}
	done

	for ML in ${LISTS}
	do
		echo " ~ adding mailinglist ${ML}"
		$CLI ml add ${TG} ${ML}

		for US in ${USERS}
		do
			echo " = adding list member ${US}"
			$CLI ml member add ${TG} ${ML} ${US}
		done
	done
done

