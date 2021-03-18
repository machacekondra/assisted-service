import utils
import deployment_options


log = utils.get_logger('deploy_squid')


def main():
    deploy_options = deployment_options.load_deployment_options()
    log.info('Staring squid deployment')
    deploy_squid(deploy_options)
    log.info('Completed squid deployment')


def deploy_squid(deploy_options):
    docs = utils.load_yaml_file_docs('deploy/squid.yaml')
    utils.set_namespace_in_yaml_docs(docs, deploy_options.namespace)
    dst_file = utils.dump_yaml_file_docs('build/squid.yaml', docs)

    log.info("Deploying %s", dst_file)
    utils.apply(
        target=deploy_options.target,
        namespace=deploy_options.namespace,
        profile=deploy_options.profile,
        file=dst_file
    )


if __name__ == "__main__":
    main()
