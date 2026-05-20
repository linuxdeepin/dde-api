// SPDX-FileCopyrightText: 2026 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#pragma once
#include <QDebug>
#include <QJsonDocument>
#include <QJsonObject>
#include <QLibrary>
#include <DSGApplication>
#include <DSysInfo>

#include <mutex>

#include <dlfcn.h>
#include <unistd.h>

namespace DDE_EventLogger {

// Check if event logging should be enabled based on system edition
// Enable for all editions except UosCommunity
inline bool shouldEnableEventLog()
{
    // Production mode: enable for all editions except UosCommunity
    // Note: DSysInfo is in Dtk::Core namespace
    return Dtk::Core::DSysInfo::uosEditionType() != Dtk::Core::DSysInfo::UosCommunity;
}

typedef bool (*Initialize)(const std::string &package_id, bool enable_sig);
typedef void (*WriteEventLog)(const std::string &event_data);

typedef struct _EventLoggerData
{
    qint64 tid;
    QString target;
    QJsonObject message;

    _EventLoggerData()
        : tid(0)
        , target(QString())
        , message(QJsonObject())
    {
    }

    _EventLoggerData(qint64 tid, const QString &target, const QJsonObject &message)
        : tid(tid)
        , target(target)
        , message(message)
    {
    }

    _EventLoggerData &operator=(const _EventLoggerData &data)
    {
        tid = data.tid;
        target = data.target;
        message = data.message;
        return *this;
    }
} EventLoggerData;

class EventLogger
{
public:
    static EventLogger &instance()
    {
        static EventLogger instance;
        return instance;
    }

    void writeEventLog(const EventLoggerData &data)
    {
        if (!shouldEnableEventLog()) {
            return;
        }

        std::lock_guard<std::mutex> lock(m_mutex);
        if (!m_initialized) {
            doInit(QString::fromUtf8(Dtk::Core::DSGApplication::id()), false);
        }
        if (!m_initialized) {
            return;
        }

        QJsonDocument json = structToJson(data);
        m_writeEventLog(json.toJson(QJsonDocument::Compact).toStdString());
    }

private:
    EventLogger()
        : m_initialized(false)
        , m_initialize(nullptr)
        , m_writeEventLog(nullptr)
        , m_library("/usr/lib/libdeepin-event-log.so")
    {
        // 检查链接库是否成功加载
        if (!m_library.load()) {
            qDebug() << "Failed to load library:" << m_library.errorString();
            return;
        }

        // 获取链接库中的函数指针
        m_initialize = (Initialize)m_library.resolve("Initialize");
        m_writeEventLog = (WriteEventLog)m_library.resolve("WriteEventLog");

        // 检查函数是否成功获取
        if (!m_initialize || !m_writeEventLog) {
            qDebug() << "Failed to resolve function";
            return;
        }
    }

    ~EventLogger() { m_library.unload(); }

    EventLogger(const EventLogger &) = delete;            // 禁止拷贝构造函数
    EventLogger &operator=(const EventLogger &) = delete; // 禁止赋值运算符

    bool doInit(const QString &package_id, bool enable_sig)
    {
        if (nullptr == m_initialize || nullptr == m_writeEventLog) {
            return false;
        }

        if (m_initialized) {
            return true;
        }

        m_initialized = m_initialize(package_id.toStdString(), enable_sig);
        if (!m_initialized) {
            qDebug() << "Failed to initialize event logger";
        }
        return m_initialized;
    }

    // 将结构体转换为 JSON 内容
    QJsonDocument structToJson(const EventLoggerData &data)
    {
        QJsonObject jsonObject;
        jsonObject["tid"] = data.tid;
        jsonObject["target"] = data.target;
        jsonObject["message"] = data.message;
        QJsonDocument jsonDocument(jsonObject);
        return jsonDocument;
    }

    std::mutex m_mutex;
    bool m_initialized;
    Initialize m_initialize;
    WriteEventLog m_writeEventLog;
    QLibrary m_library;
};
} // namespace DDE_EventLogger
