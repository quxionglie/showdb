import React, {useEffect, useState} from 'react';
import {LaptopOutlined, NotificationOutlined, UserOutlined} from '@ant-design/icons';
import {MenuProps, Table} from 'antd';
import {Breadcrumb, Layout, Menu, theme} from 'antd';
import './menu.css';

let BASE_URL = process.env.REACT_APP_BASE_URL

const {Header, Content, Footer, Sider} = Layout;
let menuItem: MenuProps['items'] = []

const App: React.FC = () => {
    const {
        token: {colorBgContainer, borderRadiusLG},
    } = theme.useToken();

    const [curTable, setCurTable] = useState({'dbName': '', 'tableName': ''});
    const [breadcrumbItems, setBreadcrumbItems] = useState([{'title': '-'},{'title': '-'}]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [dataSource, setDataSource] = useState([]);

    function fetchDatabase() {
        fetch(BASE_URL+`/api/getDatabases`)
            .then(response => {
                if (response.ok) return response.json();
                throw response;
            })
            .then(res => {
                const curMenu:any[] = []
                if (res.code == 0) {
                    for (let i = 0; i < res.data.databases.length; i++) {
                        const curDb = res.data.databases[i];
                        let curItems: any = {
                            key: curDb.dbName,
                            label: curDb.dbComment + " (" + curDb.dbType+ ")" ,
                            children: [],
                        };

                        if (curDb.tables.length > 0) {
                            for (let j = 0; j < curDb.tables.length; j++) {
                                let label = curDb.tables[j].tableName;
                                if(curDb.tables[j].tableComment){
                                    label = label+" | "+curDb.tables[j].tableComment;
                                }
                                var childrenItem = {
                                    key: curDb.dbName + "|" + curDb.tables[j].tableName,
                                    label: label,
                                };

                                curItems.children.push(childrenItem);
                            }
                        }
                        curMenu.push(curItems);
                    }
                }
                menuItem = curMenu;
            })
            .catch(err => {
                console.error(err);
                setError(err);
            })
            .finally(() => {
                setLoading(false);
            });
    }

    function fetchColumns() {
        if (curTable.tableName==''){
            return
        }

        const requestOptions = {
            method: 'GET',
            // headers: { 'Content-Type': 'application/json' },
            // body: JSON.stringify({ title: 'React POST Request Example' })
        };
        fetch(BASE_URL+`/api/getColumns?dbName=`+curTable.dbName+"&tableName="+curTable.tableName,requestOptions)
            .then(response => {
                if (response.ok) return response.json();
                throw response;
            })
            .then(res => {
                if (res.code == 0) {
                    setBreadcrumbItems([{'title': curTable.dbName},{'title': curTable.tableName}])
                    setDataSource(res.data.columns);
                }

            })
            .catch(err => {
                console.error(err);
                setError(err);
            })
            .finally(() => {
                setLoading(false);
            });
    }

    useEffect(() => {
        fetchDatabase()
    },[])
    useEffect(() => {
        fetchColumns()
    },[curTable])


    const columns = [
        {
            title: '字段名',
            dataIndex: 'columnName',
            key: 'columnName',
            width: "15%",
        },
        {
            title: '字段类型',
            dataIndex: 'columnType',
            key: 'columnType',
            width: "15%",
        },
        {
            title: '说明',
            dataIndex: 'columnComment',
            key: 'columnComment',
            width: "20%",
        },
        {
            title: '主键',
            dataIndex: 'isPk',
            key: 'isPk',
        },
        {
            title: '自增列',
            dataIndex: 'isIncrement',
            key: 'isIncrement',
        },
        {
            title: '必填',
            dataIndex: 'isRequired',
            key: 'isRequired',
        },
        {
            title: '默认值',
            dataIndex: 'columnDefault',
            key: 'columnDefault',
        },
    ];

    return (
        <Layout>
            <Header style={{display: 'flex', alignItems: 'center'}}>
                <div className="demo-logo"/>
            </Header>
            <Content style={{padding: '0 48px'}}>
                <Breadcrumb style={{margin: '16px 0'}} items={breadcrumbItems}>
                </Breadcrumb>
                <Layout
                    style={{padding: '24px 0', background: colorBgContainer, borderRadius: borderRadiusLG}}
                >
                    <Sider style={{background: colorBgContainer}} width={300}>
                        <Menu
                            mode="inline"
                            // defaultSelectedKeys={['1']}
                            // defaultOpenKeys={['sub1']}
                            style={{height: '100%'}}
                            items={menuItem}
                            onSelect={(item) => {
                                console.log("item="+item.key)
                                var items = item.key.split("|");
                                setCurTable({'dbName': items[0], 'tableName': items[1]});
                            }}
                        />
                    </Sider>
                    <Content style={{padding: '0 24px', minHeight: 280}}>
                        <Table dataSource={dataSource} columns={columns} pagination={false}/>
                    </Content>
                </Layout>
            </Content>
            <Footer style={{textAlign: 'center'}}>©2024 Created by qxl</Footer>
        </Layout>
    );
};

export default App;
