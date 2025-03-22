import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';

// Интерфейсы для типизации данных инвойса
interface Blockchain {
  name: string;
  logo: string | null;
}

interface Token {
  name: string;
  symbol: string;
  logo: string | null;
  blockchain: Blockchain;
}

interface PaymentAddress {
  public_key: string;
  requested_amount: string;
  token: Token;
}

interface Merchant {
  name: string;
}

interface InvoiceData {
  requested_amount: number;
  email: string | null;
  expires_at: string;
  merchant: Merchant;
  payment_address_list: PaymentAddress[];
}

const InvoicePage: React.FC = () => {
  // Получаем параметр payment_id из URL
  const { payment_id } = useParams<{ payment_id: string }>();

  const [invoiceData, setInvoiceData] = useState<InvoiceData | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  // Текущий шаг: 1 - выбор валюты, 2 - подтверждение суммы, 3 - вывод адреса кошелька
  const [currentStep, setCurrentStep] = useState<number>(1);
  const [selectedPayment, setSelectedPayment] = useState<PaymentAddress | null>(null);
  const [timeLeft, setTimeLeft] = useState<number>(0);

  // Запрос данных инвойса с сервера
  useEffect(() => {
    const fetchInvoice = async () => {
      try {
        const response = await fetch(`http://127.0.0.1:8080/frontend/api/v1/payments/${payment_id}/`);
        if (!response.ok) {
          throw new Error('Ошибка загрузки данных инвойса');
        }
        const data: InvoiceData = await response.json();
        setInvoiceData(data);
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchInvoice();
  }, [payment_id]);

  // Таймер до expires_at (если задан)
  useEffect(() => {
    if (invoiceData && invoiceData.expires_at) {
      const interval = setInterval(() => {
        const expires = new Date(invoiceData.expires_at);
        const now = new Date();
        const diff = Math.max(0, Math.floor((expires.getTime() - now.getTime()) / 1000));
        setTimeLeft(diff);
      }, 1000);
      return () => clearInterval(interval);
    }
  }, [invoiceData]);

  const formatTime = (seconds: number) => {
    const m = Math.floor(seconds / 60);
    const s = seconds % 60;
    return `${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`;
  };

  // Обработчик выбора валюты
  const handleCurrencyChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const index = e.target.value;
    if (invoiceData && index !== '') {
      const payment = invoiceData.payment_address_list[parseInt(index)];
      setSelectedPayment(payment);
      setCurrentStep(2);
    }
  };

  // Переход к шагу 3 по нажатию на кнопку "Оплатить"
  const handlePayClick = () => {
    setCurrentStep(3);
  };

  // Обработчик копирования адреса в буфер обмена
  const handleCopyAddress = () => {
    if (selectedPayment) {
      navigator.clipboard.writeText(selectedPayment.public_key)
        .then(() => alert('Адрес скопирован'))
        .catch(() => alert('Ошибка копирования'));
    }
  };

  if (loading) {
    return <div style={outerContainerStyle}>Загрузка...</div>;
  }

  if (error) {
    return <div style={outerContainerStyle}>Ошибка: {error}</div>;
  }

  if (!invoiceData) {
    return <div style={outerContainerStyle}>Нет данных инвойса</div>;
  }

  return (
    <div style={outerContainerStyle}>
      <div style={cardStyle}>
        {/* Логотип сайта */}
        <div style={logoContainerStyle}>
          <img src="https://via.placeholder.com/150?text=LOGO" alt="Logo" style={{ maxWidth: '100%' }} />
        </div>

        {/* Таймер, если задан expires_at */}
        {invoiceData.expires_at && (
          <div style={timerStyle}>
            Время до истечения: {formatTime(timeLeft)}
          </div>
        )}

        {/* Вывод основной суммы */}
        <div style={amountStyle}>
          Сумма: {invoiceData.requested_amount}
        </div>

        {/* Шаг 1: Выбор валюты */}
        {currentStep === 1 && (
          <div style={{ marginBottom: '20px' }}>
            <select onChange={handleCurrencyChange} style={selectStyle}>
              <option value="">Выбор валюты</option>
              {invoiceData.payment_address_list.map((payment, index) => (
                <option key={index} value={index}>
                  {payment.token.name} ({payment.token.symbol})
                </option>
              ))}
            </select>
          </div>
        )}

        {/* Шаг 2: Вывод суммы для выбранной валюты и разблокировка кнопки "Оплатить" */}
        {currentStep === 2 && selectedPayment && (
          <div style={{ marginBottom: '20px' }}>
            <div style={amountStyle}>
              {selectedPayment.requested_amount} {selectedPayment.token.symbol}
            </div>
            <button onClick={handlePayClick} style={buttonEnabledStyle}>
              Оплатить
            </button>
          </div>
        )}

        {/* Шаг 3: Вывод адреса кошелька и кнопки копирования */}
        {currentStep === 3 && selectedPayment && (
          <div>
            <div style={infoTextStyle}>Отправьте оплату на следующий адрес:</div>
            <div style={addressStyle}>{selectedPayment.public_key}</div>
            <button onClick={handleCopyAddress} style={buttonEnabledStyle}>
              Скопировать адрес
            </button>
          </div>
        )}

        {/* На шаге 1 кнопка "Оплатить" заблокирована */}
        {currentStep === 1 && (
          <div>
            <button disabled style={buttonDisabledStyle}>
              Оплатить
            </button>
          </div>
        )}
      </div>
    </div>
  );
};

// Стили

const outerContainerStyle: React.CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  minHeight: '100vh',
  backgroundColor: '#000',
  padding: '20px',
  color: '#fff'
};

const cardStyle: React.CSSProperties = {
  backgroundColor: '#111',
  padding: '30px',
  borderRadius: '8px',
  boxShadow: '0 4px 8px rgba(0,0,0,0.3)',
  maxWidth: '400px',
  width: '100%',
  textAlign: 'center',
  color: '#fff'
};

const logoContainerStyle: React.CSSProperties = {
  marginBottom: '20px'
};

const timerStyle: React.CSSProperties = {
  marginBottom: '15px',
  fontSize: '14px',
  color: '#aaa'
};

const amountStyle: React.CSSProperties = {
  fontSize: '24px',
  marginBottom: '20px'
};

const selectStyle: React.CSSProperties = {
  padding: '10px',
  borderRadius: '4px',
  border: 'none',
  width: '100%',
  fontSize: '16px',
  backgroundColor: '#222',
  color: '#fff'
};

const buttonEnabledStyle: React.CSSProperties = {
  padding: '10px 20px',
  backgroundColor: '#222',
  color: '#fff',
  border: '1px solid #fff',
  borderRadius: '4px',
  cursor: 'pointer',
  fontSize: '16px',
  width: '100%'
};

const buttonDisabledStyle: React.CSSProperties = {
  ...buttonEnabledStyle,
  opacity: 0.5,
  cursor: 'not-allowed'
};

const addressStyle: React.CSSProperties = {
  backgroundColor: '#000',
  border: '1px solid #fff',
  padding: '10px',
  borderRadius: '4px',
  marginBottom: '20px',
  wordBreak: 'break-all',
  color: '#fff'
};

const infoTextStyle: React.CSSProperties = {
  marginBottom: '10px'
};

export default InvoicePage;