import { ConfigProvider, Layout } from 'antd';
import zhCN from 'antd/locale/zh_CN';

function App() {
    return <ConfigProvider locale={zhCN}
    theme={{
        components: {
            Layout: {
                headerBg: "#f0f0f0",
                headerHeight: 44,
                headerPadding: "0 8px 0 88px",
                siderBg: "#f0f0f0",
            }
        }
    }}
    >
        <Layout style={{
            height: "100vh"
        }}>
            <Layout.Header style={{borderBottom: "1px solid #e0e0e0"}}>
                hello navidog
            </Layout.Header>
            <Layout>
                <Layout.Sider></Layout.Sider>
                <Layout.Content></Layout.Content>
            </Layout>
        </Layout>
    </ConfigProvider>
}

export default App
