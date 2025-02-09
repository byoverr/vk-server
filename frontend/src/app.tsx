import React, { useEffect, useState } from "react";
import { Table, notification } from "antd";
import type { ColumnsType } from "antd/es/table";

interface Container {
  id: number;
  ip: string;
  ping_time: string;
  last_check: string;
}

const App: React.FC = () => {
  const [containers, setContainers] = useState<Container[]>([]);

  // Mock данные для заглушки
  const mockData: Container[] = [
    {
      id: 1,
      ip: "192.168.1.1",
      ping_time: "1.23ms",
      last_check: "2023-10-01T12:34:56Z",
    },
    {
      id: 2,
      ip: "192.168.1.2",
      ping_time: "2.34ms",
      last_check: "2023-10-01T12:35:00Z",
    },
  ];

  // Функция для загрузки данных
  const fetchData = async () => {
    try {
      const response = await fetch("/containers");
      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      const data = await response.json();
      setContainers(data);
    } catch (error) {
      console.error("Error fetching data:", error);
      notification.error({
        message: "Ошибка",
        description: "Не удалось загрузить данные. Используются mock-данные.",
      });
      setContainers(mockData); // Использовать заглушку при ошибке
    }
  };

  // Загрузка данных при монтировании и каждые 60 секунд
  useEffect(() => {
    fetchData(); // Первый запрос
    const interval = setInterval(fetchData, 60000); // Периодический запрос

    return () => clearInterval(interval); // Очистка интервала при размонтировании
  }, []);

  // Определение колонок таблицы
  const columns: ColumnsType<Container> = [
    {
      title: "IP",
      dataIndex: "ip",
      key: "ip",
    },
    {
      title: "Ping Time",
      dataIndex: "ping_time",
      key: "ping_time",
    },
    {
      title: "Last Check",
      dataIndex: "last_check",
      key: "last_check",
    },
  ];

  return (
    <div style={{ padding: "20px" }}>
      <Table
        dataSource={containers}
        columns={columns}
        rowKey="id"
        bordered
        pagination={{ pageSize: 10 }}
      />
    </div>
  );
};

export default App;