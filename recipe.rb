HOST = node[:hostname]
USER = 'isucon'

service 'isucari.golang.service' do
	action %i[enable start]
end

###
# Monitoring tools
###

execute 'install netdata' do
	command 'bash -c "bash <(curl -Ss https://my-netdata.io/kickstart.sh) all --dont-wait"'
	not_if "systemctl list-unit-files | grep '^netdata.service'"
end

service 'netdata' do
        action %i[disable stop]
end

package 'percona-toolkit'

###
# Configure middlewares
###

#package 'redis'

#service 'redis' do
#	action %i[enable start]
#end

if HOST == 'isu01' || HOST == 'isu02'
	service 'mysql' do
		action %i[disable stop]
	end
end

{
	'nginx' => '/etc/nginx/nginx.conf',
	'mysql' => '/etc/mysql/mysql.conf.d/mysqld.cnf',
}.each do |service, conf|
	conf_source = "./#{HOST}/#{File.basename conf}"
	next unless File.file? conf_source

	service service do
		action %i[enable start]
	end

	remote_file conf do
		source conf_source
		mode '0644'
		notifies :restart, "service[#{service}]"
	end
end
