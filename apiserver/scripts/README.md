# Digital Ocean deployment scripts

* create_load_balancer.sh
Creates a new load balancer
  * Param: 
    * 1: Name of new load balancer
    
* create_vm.sh
Creates a new VM, starts a docker container, and adds to the hardcoded load balancer
  * Params:   
    * 1: VM name 
    * 2: Docker container name
    * 3: Docker image name

* add_droplet_to_load_balancer.sh
Adds droplet to hardcoded load balancer
  * Param: 
    * 1: Droplet id's (CSV)
    
* get_droplets.sh
Returns all droplets for current API token
  
* get_load_balancer.sh
Returns the requested load balancer
   * Param:
     * 1: Load balancer id
     
* remove_droplets_from_load_balancer.sh
Gets first load balancer for current API token and prompts for removal/destruction of droplets from it

