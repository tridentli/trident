#!/bin/sh

### BEGIN INIT INFO
# Provides:          trident
# Required-Start:    $remote_fs $syslog $network $time
# Required-Stop:     $remote_fs $syslog $network $time
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: Trident - Trusted Information Exchange Toolkit
# Description:
#   The Trident Trusted Information Exchange Toolkit
#   Please see https://trident.li
### END INIT INFO

set -e

DESC="Trident - Trusted Information Exchange Toolkit"
NAME=trident
DNAME=tridentd
DAEMON=/usr/sbin/${DNAME}
DAEMON_OPTS=""
PIDDIR=/var/run/trident
PIDFILE=${PIDDIR}/${DNAME}.pid
SCRIPTNAME=/etc/init.d/${NAME}
CONFFILE=/etc/trident/${NAME}.conf
DAEMON_USER=${NAME}

# Exit if the package is not installed
[ -x "${DAEMON}" ] || exit 0

# Read configuration variable file if it is present
[ -r /etc/default/${DNAME} ] && . /etc/default/${DNAME}

# Add daemonize && username & pidfile options
DAEMON_OPTS="${DAEMON_OPTS} --daemonize --username ${DAEMON_USER} --pidfile ${PIDFILE}"

# Define LSB log_* functions.
# Depend on lsb-base (>= 3.0-6) to ensure that this file is present.
. /lib/lsb/init-functions

# Make sure the PID dir is there
mkdir -p ${PIDDIR}
chown trident:trident ${PIDDIR}
chmod 755 ${PIDDIR}

check_for_upstart() {
    if init_is_upstart; then
        exit $1
    fi
}

check_for_no_start() {
	# Is Trident enabled?
	case "${TRIDENT_ENABLED}" in
	[Nn]*)
		exit 0
		;;
	esac
}

check_config() {
	if [ ! -f ${CONFFILE} ];
	then
		log_failure_msg "Trident Configuration file ${CONFFILE} doesn't exist" || true
		exit 1;
	fi
}

case "$1" in
  start)
	check_for_upstart 1
	check_for_no_start
	log_daemon_msg "Starting ${DESC}" "${NAME}" || true
	if start-stop-daemon --start --quiet --oknodo --pidfile ${PIDFILE} --exec ${DAEMON} -- ${DAEMON_OPTS};
	then
		log_end_msg 0 || true
	else
		log_end_msg 1 || true
        fi
	;;

  stop)
	check_for_upstart 0
	log_daemon_msg "Stopping ${DESC}" "${DNAME}" || true
	if start-stop-daemon --stop --quiet --oknodo --pidfile ${PIDFILE};
	then
		log_end_msg 0 || true
	else
		log_end_msg 1 || true
	fi
	;;

  reload|force-reload)
	log_action_msg "Reloading Trident is not possible" || true
	;;

  restart)
	check_for_upstart 1
	check_config
	log_daemon_msg "Restarting ${DESC}" "${DNAME}"
	start-stop-daemon --stop --quiet --oknodo --retry 30 --pidfile ${PIDFILE}
	check_for_no_start log_end_msg
	if start-stop-daemon --start --quiet --oknodo --pidfile ${PIDFILE} --exec ${DAEMON} -- ${DAEMON_OPTS}; then
		log_end_msg 0 || true
	else
		log_end_msg 1 || true
	fi
	;;

  status)
	check_for_upstart 1
	status_of_proc "${DAEMON}" "${DNAME}" && exit 0 || exit $?
	;;

  *)
	log_action_msg "Usage: ${SCRIPTNAME} {start|stop|restart|status}" || true
	exit 1
	;;
esac

