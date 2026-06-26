import { Link, useNavigate } from 'react-router-dom'

// Общая навигационная панель страниц с .navbar — заменяет дублированную разметку,
// которая раньше повторялась почти одинаково в каждой странице.
//
// Props:
//   onBack    — если передан, показывает кнопку "← Назад" с этим обработчиком
//   backLabel — текст кнопки назад (по умолчанию "← Назад")
//   title     — текст бренда/заголовка слева
//   titleTo   — если задан, title рендерится как Link на этот путь
//   links     — [{ to, label, badge }] — обычные нав-ссылки
//   right     — произвольный JSX в конце навбара (кнопки, колокольчик и т.п.)
export default function Navbar({ onBack, backLabel = '← Назад', title, titleTo, links = [], right }) {
  const navigate = useNavigate()

  return (
    <nav className="navbar">
      {onBack && (
        <button onClick={onBack === true ? () => navigate(-1) : onBack} className="back-btn">
          {backLabel}
        </button>
      )}
      {title && (
        titleTo
          ? <Link to={titleTo} className="navbar-brand">{title}</Link>
          : <span className="navbar-brand">{title}</span>
      )}
      {(links.length > 0 || right) && (
        <div className="navbar-links">
          {links.map(link => (
            <Link key={link.to} to={link.to} className={link.badge !== undefined ? 'navbar-badge-link' : undefined}>
              {link.label}
              {link.badge > 0 && <span className="navbar-badge">{link.badge}</span>}
            </Link>
          ))}
          {right}
        </div>
      )}
    </nav>
  )
}
