import { CheckCircle2Icon } from "lucide-react"

import {
  Alert,
  AlertTitle,
} from "@/components/ui/alert"
import "./alert.css"

export const CustomAlert = ({alertText} : {alertText: string}) => {
  return (
    <div className="custom-alert">
      <Alert>
        <CheckCircle2Icon />
        <AlertTitle>{alertText}</AlertTitle>
      </Alert>
    </div>
  )
}
