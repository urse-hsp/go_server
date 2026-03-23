import { LockOutlined, UserOutlined } from '@ant-design/icons'
import {
  LoginFormPage,
  ProConfigProvider,
  ProFormText,
} from '@ant-design/pro-components'
import { Tabs, message, theme } from 'antd'
import { useState } from 'react'
import { useRef } from 'react'
import type { ProFormInstance } from '@ant-design/pro-components'

type LoginType = 'create' | 'account'

const Page = () => {
  const formRef = useRef<ProFormInstance>()

  const [loginType, _setLoginType] = useState<LoginType>('account')
  const { token } = theme.useToken()

  return (
    <div
      style={{
        backgroundColor: 'white',
        height: '100vh',
      }}
    >
      <LoginFormPage
        formRef={formRef}
        backgroundVideoUrl='https://gw.alipayobjects.com/v/huamei_gcee1x/afts/video/jXRBRK_VAwoAAAAAAAAAAAAAK4eUAQBr'
        title={loginType === 'account' ? '登录' : '注册'}
        containerStyle={{
          backgroundColor: 'rgba(0, 0, 0,0.65)',
          backdropFilter: 'blur(4px)',
        }}
        onFinish={async (values) => {
          console.log('表单数据:', values) // 👈 这里就能拿到数据！
          // 后续在这里调用 API
        }}
      >
        <Tabs
          items={[
            { key: 'account', label: '账号密码登录' },
            { key: 'create', label: '注册' },
          ]}
          onChange={(activeKey) => {
            _setLoginType(activeKey as LoginType)
            formRef.current?.resetFields?.()
          }}
        />
        {(loginType === 'account' || loginType === 'create') && (
          <>
            <ProFormText
              name='username'
              fieldProps={{
                size: 'large',
                prefix: (
                  <UserOutlined
                    style={{
                      color: token.colorText,
                    }}
                    className={'prefixIcon'}
                  />
                ),
              }}
              placeholder={'用户名: admin or user'}
              rules={[
                {
                  required: true,
                  message: '请输入用户名!',
                },
              ]}
            />
            <ProFormText.Password
              name='password'
              fieldProps={{
                size: 'large',
                prefix: (
                  <LockOutlined
                    style={{
                      color: token.colorText,
                    }}
                    className={'prefixIcon'}
                  />
                ),
              }}
              placeholder={'密码: ant.design'}
              rules={[
                {
                  required: true,
                  message: '请输入密码！',
                },
              ]}
            />
          </>
        )}

        {/* {loginType === 'phone' && (
          <>
            <ProFormText
              fieldProps={{
                size: 'large',
                prefix: (
                  <MobileOutlined
                    style={{
                      color: token.colorText,
                    }}
                    className={'prefixIcon'}
                  />
                ),
              }}
              name="mobile"
              placeholder={'手机号'}
              rules={[
                {
                  required: true,
                  message: '请输入手机号！',
                },
                {
                  pattern: /^1\d{10}$/,
                  message: '手机号格式错误！',
                },
              ]}
            />
            <ProFormCaptcha
              fieldProps={{
                size: 'large',
                prefix: (
                  <LockOutlined
                    style={{
                      color: token.colorText,
                    }}
                    className={'prefixIcon'}
                  />
                ),
              }}
              captchaProps={{
                size: 'large',
              }}
              placeholder={'请输入验证码'}
              captchaTextRender={(timing, count) => {
                if (timing) {
                  return `${count} ${'获取验证码'}`;
                }
                return '获取验证码';
              }}
              name="captcha"
              rules={[
                {
                  required: true,
                  message: '请输入验证码！',
                },
              ]}
              onGetCaptcha={async () => {
                message.success('获取验证码成功！验证码为：1234');
              }}
            />
          </>
        )} */}
        {/* <div
          style={{
            marginBlockEnd: 24,
          }}
        >
          <ProFormCheckbox noStyle name="autoLogin">
            自动登录
          </ProFormCheckbox>
          <a
            style={{
              float: 'right',
            }}
          >
            忘记密码
          </a>
        </div> */}
      </LoginFormPage>
    </div>
  )
}

const Demo = () => {
  return (
    <ProConfigProvider dark>
      <Page />
    </ProConfigProvider>
  )
}

export default () => (
  <div style={{ padding: 24 }}>
    <Demo />
  </div>
)
