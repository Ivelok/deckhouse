{% if page.platform_code == 'openstack' or page.platform_code == 'vsphere' %}
{% assign ee_only = true %}
{% else %}
{% assign ee_only = false %}
{% endif %}
## Схема установки

Установка Deckhouse Platform {% if page.platform_type == 'baremetal' %}на{% else %}в{% endif %} {{ page.platform_name }} в общем случае выглядит так:
-  На локальной машине (с которой будет производиться установка) запускается Docker-контейнер.
-  Этому контейнеру передаются приватный SSH-ключ с локальной машины и файл конфигурации будущего кластера в формате YAML (например, `config.yml`).
-  Контейнер подключается по SSH к целевой машине (для bare metal-инсталляций) или облаку, после чего происходит непосредственно установка и настройка кластера Kubernetes.
{% if page.platform_type == 'cloud' %}
> При установке Deckhouse в публичное облако для Kubernetes-кластера будут использоваться «обычные» вычислительные ресурсы провайдера, а не managed-решение с Kubernetes, предлагаемое провайдером.
{%- endif %}
>
> Подробнее о каналах обновления Deckhouse Platform (release channels) можно почитать в [документации](/ru/documentation/v1/deckhouse-release-channels.html).

## Ограничения/требования для установки:

-   На машине, с которой будет производиться установка, необходимо наличие Docker runtime.
-   Доступ к интернет и стандартным репозиториям используемой ОС для установки дополнительных необходимых пакетов.
-   Deckhouse поддерживает разные версии Kubernetes: с 1.16 по 1.21 включительно. Однако обратите внимание, что для установки «с чистого листа» доступны только версии 1.16, 1.19, 1.20 и 1.21. В примерах конфигурации ниже используется версия 1.19.
{% if page.platform_type == 'cloud' %}
-   Рекомендуется использовать типы инстанса master-узлов для будущего кластера с характеристиками не хуже следующих:
{%- else %}
-   Рекомендованная минимальная аппаратная конфигурация master-узлов для будущего кластера:
{%- endif %}
    -   не менее 4 ядер CPU;
    -   не менее 8  ГБ RAM;
    -   не менее 40 ГБ дискового пространства;
    -   ОС: Ubuntu Linux 16.04/18.04/20.04 LTS или CentOS 7.

{% if ee_only != true %}
## Выберите редакцию Deckhouse Platform для продолжения установки в {{ page.platform_name }}

[Сравнить](/ru/products/enterprise_edition.html#ce-vs-ee) возможности редакций Community Edition и Enterprise Edition.
{% endif %}