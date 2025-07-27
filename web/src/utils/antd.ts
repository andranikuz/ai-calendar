// Centralized Ant Design imports for better tree shaking
// Using ES modules for optimal tree shaking with Vite

// Re-export individual components from their ES modules
export { default as Alert } from 'antd/es/alert';
export { default as App } from 'antd/es/app';
export { default as Avatar } from 'antd/es/avatar';
export { default as Badge } from 'antd/es/badge';
export { default as Button } from 'antd/es/button';
export { default as Calendar } from 'antd/es/calendar';
export { default as Card } from 'antd/es/card';
export { default as Checkbox } from 'antd/es/checkbox';
export { default as Col } from 'antd/es/col';
export { default as ConfigProvider } from 'antd/es/config-provider';
export { default as DatePicker } from 'antd/es/date-picker';
export { default as Divider } from 'antd/es/divider';
export { default as Drawer } from 'antd/es/drawer';
export { default as Dropdown } from 'antd/es/dropdown';
export { default as Empty } from 'antd/es/empty';
export { default as Form } from 'antd/es/form';
export { default as Grid } from 'antd/es/grid';
export { default as Input } from 'antd/es/input';
export { default as InputNumber } from 'antd/es/input-number';
export { default as Layout } from 'antd/es/layout';
export { default as List } from 'antd/es/list';
export { default as Menu } from 'antd/es/menu';
export { default as Modal } from 'antd/es/modal';
export { default as Popconfirm } from 'antd/es/popconfirm';
export { default as Progress } from 'antd/es/progress';
export { default as Row } from 'antd/es/row';
export { default as Select } from 'antd/es/select';
export { default as Slider } from 'antd/es/slider';
export { default as Space } from 'antd/es/space';
export { default as Spin } from 'antd/es/spin';
export { default as Statistic } from 'antd/es/statistic';
export { default as Switch } from 'antd/es/switch';
export { default as Tag } from 'antd/es/tag';
export { default as TimePicker } from 'antd/es/time-picker';
export { default as Tooltip } from 'antd/es/tooltip';
export { default as Tree } from 'antd/es/tree';
export { default as Typography } from 'antd/es/typography';
export { default as message } from 'antd/es/message';
export { default as notification } from 'antd/es/notification';

// Export types
export type {
  FormInstance,
  SelectProps,
  ModalProps,
  CardProps,
  ButtonProps,
} from 'antd';