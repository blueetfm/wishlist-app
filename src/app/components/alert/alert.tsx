import { CheckCircle2Icon } from "lucide-react"

import {
  Alert,
  AlertTitle,
} from "@/app/components/ui/alert"

export const CustomAlert = ({alertText} : {alertText: string}) => {
  return (
    <div className="grid w-full max-w-xl items-start gap-4">
      <Alert>
        <CheckCircle2Icon />
        <AlertTitle>{alertText}</AlertTitle>
      </Alert>
    </div>
  )
}
