# Digital Ocean deployment scripts

* create_load_balancer.sh
  * Param: 
    * 1: Name of new load balancer
    
* create_vm.sh
  * Params:   
    * 1: VM name 
    * 2: Docker container name
    * 3: Docker image name

* add_droplet_to_load_balancer.sh
  * Param: 
    * 1: Droplet id's (CSV)
    
* get_droplets.sh
  * Returns all droplets for current API token
  
* get_load_balancer.sh
   * Param:
     * 1: Load balancer id
     
* remove_droplets_from_load_balancer.sh

